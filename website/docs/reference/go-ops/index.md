---
title: Go ops
sidebar_label: Overview
---

A Go op is a directory that contains `op.go` in `package main`, with a `go.mod` in that directory or a parent. It can't also contain `op.yml`. The program
imports the dependency-free Go op SDK, `github.com/opctl/opctl/sdks/go/opspec/sdk`, conventionally as `opctl`. For a walkthrough, see [Authoring ops in
Go](../../training/authoring-ops-in-go.md).

- [Inputs and outputs](inputs-and-outputs.md): declaring an op's interface.
- [Op stubs](op-stubs.md): typed op calls to YAML and Go ops.
- [Calls](calls.md): container calls, op calls, and groups.
- [Handles](handles.md): passing values between inputs, calls, and outputs.
- [Bootstrap](bootstrap.md): how opctl runs the program, and its caches.

## Entry points

- `Main[I any](fn func(I) Call, op ...Op)`
- `MainWithOutputs[I, O any](fn func(I) (Call, O), op ...Op)`

`main` calls one of these with the op's plan function. `I` is the [inputs struct](inputs-and-outputs.md#inputs); an op without inputs uses `struct{}`.
`MainWithOutputs` also returns a struct of [output handles](inputs-and-outputs.md#outputs). The program writes only its response to stdout, and diagnostics to
stderr.

## Name and description

The op's name is its directory name, and its description is the package doc comment, rendered as markdown. Build constraints apply, so a file tagged `//go:build
ignore` doesn't contribute its comment. Test files never do.

`opctl.Op` overrides either one. Pass it as the last argument to `Main`, `MainWithOutputs`, `Describe`, or `DescribeWithOutputs`:

```go
func main() {
    opctl.Main(func(struct{}) opctl.Call {
        return opctl.Container("alpine:3").Cmd("echo", "deploying")
    }, opctl.Op{Name: "deploy", Description: "Deploys the API."})
}
```

Empty fields keep the default. Passing more than one `Op` is an error. Go ops have no `version` field; version them with git tags.

## Describe and plan

opctl runs the program in two phases, together called [bootstrap](bootstrap.md):

- **Describe** obtains the op's interface: its name, description, inputs, and outputs. It receives no input values and never calls the plan function.
- **Plan** calls the plan function with one call's resolved input values, and obtains the op graph that opctl then runs.

The plan function only builds the op graph; building it runs nothing. It runs again for every call, so it must build the same graph for the same inputs. Don't
depend on map iteration order, time, randomness, or goroutines.

## Testing without a node

- `Plan[I any](fn func(I) Call, inputs map[string]any) ([]byte, error)`
- `PlanWithOutputs[I, O any](fn func(I) (Call, O), inputs map[string]any) ([]byte, error)`
- `Describe[I any](op ...Op) ([]byte, error)`
- `DescribeWithOutputs[I, O any](op ...Op) ([]byte, error)`

These return the same JSON as the plan and describe phases, in process. `Plan` takes already resolved values: pass every public input, including those with
defaults. Secret, directory, file, and socket inputs stay handles and need no value. Both use the working directory's name as the op name, and `Describe` reads
the doc comment of the package in the working directory, so `go test` in the op directory matches what opctl sees.

## Using Go ops

Run, validate, and list Go ops with the same commands as YAML ops:

```sh
opctl run -a env=prod -a token=example path/to/op
opctl op validate path/to/op
opctl ls path/to/op
```

- `opctl run` prompts for missing inputs, applies defaults, and checks constraints. An invalid explicit value prompts again, or fails in a non-interactive
  terminal; it never falls back to the default. Secret prompts are masked.
- `opctl op validate` checks the described interface without calling the plan function.
- `opctl ls` lists Go ops with their descriptions.
- The UI shows a Go op's name and description when you open its `op.go`. Its op graph appears when the op runs.

`opctl op create --go [-p PATH] [-d DESCRIPTION] NAME` creates `NAME/op.go`, `NAME/go.mod`, and a `.gitignore` containing `/opctlgen/`, with one input and a container call.
The description becomes the
package doc comment, and `NAME` must be a valid Go module path. The generated module path is `ops.local/NAME`, so op names such as `fmt` don't collide with the
standard library. The command writes nothing if the directory already contains `op.yml`, `op.go`, `go.mod`, `go.sum`, or `.gitignore`. A release build of opctl requires its
matching Go op SDK version in the scaffold and writes its `go.sum`. Otherwise it requires version `v0.0.0`, and opctl prints a reminder to add a `replace`
directive for `github.com/opctl/opctl/sdks/go/opspec/sdk` that points it at `sdks/go/opspec/sdk` in an opctl checkout.

## Schema negotiation

Describe and plan each receive the range of opspec schema versions the node can read. The Go op SDK supports schemas 1 through 4, and emits the newest version
both support. A graph that uses only baseline features emits schema 1 for a node that reads `[1,1]` and schema 4 for a node that reads `[1,4]`, with the same
calls, mounts, and output bindings.

| Required schema | Graph features |
| --- | --- |
| 1 | Calls, groups, mounts, op inputs and outputs |
| 2 | Container stdout and stderr bindings |
| 3 | Outputs file members |
| 4 | Typed container outputs |

The SDK inspects the graph built for the resolved inputs, so unused builders don't require a newer schema. It never drops a feature to fit an older node. A
graph that needs a newer schema than the node reads fails before any of its calls start, with an error that names the feature and its location in the graph:

```text
this op uses typed container outputs (opspec >= 4); upgrade opctl to a version supporting opspec 4 (at run.serial[0].container.outputParams)
```

Outputs file members on an older node report `this op uses container outputs file (opspec >= 3)`, and stdout or stderr bindings report `this op uses container
stdout/stderr bindings (opspec >= 2)`. If the node requires a newer schema than the SDK supports, the error names both ranges and says `upgrade the Go op SDK`.
Schema versions are integers, independent of the YAML op file's `opspec` version.

Describe checks the interface on its own, so `opctl op validate` can succeed while `opctl run` rejects an input-dependent graph that needs a newer schema. The
in-process `Plan` and `Describe` helpers use the SDK's whole range.
