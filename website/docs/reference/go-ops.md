---
title: Go ops
---

A Go op is a directory containing `op.go` and a `go.mod` in that directory or a parent. It cannot also contain `op.yml`. Import the dependency-free module `github.com/opctl/opctl/sdks/go/opspec/sdk` to build an op graph.

```go
package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
    Env opctl.String `opctl:"env,description='Target environment',default=dev,enum=dev|prod"`
    Token opctl.String `opctl:"token,secret,minLength=1"`
}

func main() {
    opctl.Main(func(in Inputs) opctl.Call {
        return opctl.Container("alpine:3").Cmd("echo", in.Env.Value()).
            EnvFrom("TOKEN", in.Token)
    })
}
```

`String`, `Number`, and `Boolean` handles expose their resolved input through `Value()`. `Array[T].Value()` returns `[]T`; `Object[T].Value()` returns the struct or string-keyed map `T`. The SDK derives item and property schemas from Go types and uses `json` tags for property names. `[]any` and `map[string]any` allow untyped values. Recursive types, pointers, and embedded struct fields are rejected.

Every input field must be exported, use a handle type, and have a unique `opctl` tag name containing letters, numbers, `-`, or `_`. The `opctl_` prefix is reserved for generated handles. Tags accept `description`, `default`, and `secret`. Boolean inputs cannot be secret; sockets cannot have defaults. Single-quoted tag values can contain commas and escape single quotes or backslashes.

String constraints include `enum`, `format`, `pattern`, `minLength`, and `maxLength`. Number constraints include `enum`, `format=integer`, `minimum`, `maximum`, and `multipleOf`. Arrays accept `minItems`, `maxItems`, and `uniqueItems`. Objects accept `minProperties`, `maxProperties`, and `required=a|b`. Use `constraints='<JSON object>'` for nested JSON Schema constraints. Duplicate keys and unknown tag keys are errors.

Use `EnvFrom` to bind handles, `Dir` and `File` to mount inputs, and `Socket` to bind a socket. Secret input values cannot be read during plan. Directory, file, and socket handles have no `Value()` method. Bootstrap receives only non-secret scalar and structured values; filesystem inputs remain handles. Literal strings escape `$` so input text cannot introduce scope references.

Unsigned Go integer types emit `minimum: 0` in array item and object property schemas.

The `constraints` tag can add JSON schema keywords but cannot replace keywords derived from the Go input type.

`Describe[Inputs]()` returns the interface without invoking the plan function. `Plan(fn, inputs)` builds the graph in process with already resolved public input values. No-input ops use `func(struct{}) opctl.Call`. Use `Serial(calls...)` to run calls in order and `Parallel(calls...)` to run them concurrently. Helpers are plain Go functions returning calls. Go loops can build a slice of calls to pass to either group. Empty groups are dropped from their parent; an empty root graph is an error.

Run a Go op with the same arguments as a YAML op:

```sh
opctl run -a env=prod -a token=example path/to/op
opctl op validate path/to/op
opctl ls path/to/op
```

The CLI prompts for missing inputs, applies defaults, and checks constraints. Invalid explicit Go input values prompt again or fail in a non-interactive terminal; they do not fall back to defaults. Secret prompts are masked. Directory and file defaults such as `$(./testdata)` resolve relative to the op directory. `opctl op validate` checks the declared interface without invoking the plan function. `opctl ls` includes Go ops and their descriptions. Opening `op.go` in the webapp shows its described name and description; its run graph is built when the op runs.

The node bootstraps describe and plan in containers using the Go version in the nearest `go.mod`. The module root and each local `replace` target mount separately at their host paths, including absolute replacements. Mounts that hide required bootstrap image directories are rejected. Bootstrap calls appear beneath the op call. Bootstrap stdout is captured in memory, with Docker logging disabled, and is never published as log events. Bootstrap failures return up to 64 KiB of stderr diagnostics to the client. Describe is also available through `POST /api/ops/describes` with `{"op":{"ref":"path/to/op"}}`. Invalid declarations, bootstrap program exits, and response parse failures return HTTP 400; runtime and context failures return HTTP 500.

Describe and plan compile separately for each request. No compiled-program cache is available yet. Go ops do not yet declare op outputs.

```go
func checks() opctl.Call {
    var calls []opctl.Call
    for _, name := range []string{"lint", "test"} {
        calls = append(calls, opctl.Container("alpine:3").Cmd("echo", name))
    }
    return opctl.Serial(opctl.Parallel(calls...), opctl.Container("alpine:3").Cmd("echo", "done"))
}
```

`NewDir()` and `NewFile()` create output handles inside the function passed to `Main` or `Plan`. The first container to mount a handle produces it; later serial calls mount the same directory or file. Parallel children start with the same incoming scope, and their outputs become visible after the group. Plan rejects a handle produced by parallel siblings, or read by `EnvFrom` before an earlier call mounts it.

```go
out := opctl.NewDir()
build := opctl.Container("alpine:3").Cmd("sh", "-ce", "echo hello > /out/result.txt").Dir("/out", out)
check := opctl.Container("alpine:3").Cmd("cat", "/in/result.txt").Dir("/in", out)
return opctl.Serial(build, check)
```

Mounts are shared and mutable. Parallel consumers may mount an already visible handle, but their writes need coordination. Prefer directory outputs over file outputs: tools that replace a file by renaming a temporary file cannot replace a bind-mounted file. Writing a mounted file in place, for example with shell `>`, works.

Build the graph deterministically in one goroutine, and sort map keys before iterating over them. Generated names are stable for repeated plans with the same handle creation order. Create handles inside the plan function, including inside helpers and Go loops.

Use `RawRef[opctl.File]("$(./expected.txt)")` or `RawRef[opctl.Dir]("$(./scripts)")` for op-directory paths. `RawRef[T]` accepts only a whole `$(...)` reference and emits its inner text verbatim, including nested references. It works with any handle type, but bypasses visibility checks and has no readable plan-time value. Building a raw reference from input text can inject references into the graph. Prefer typed handles when they can express the binding. Literal commands, environment values, working directories, and image strings continue to escape `$`.

## Calling ops by ref

`CallOp("./yaml").Inputs(map[string]any{"label": in.Label})` builds an op call.
Refs resolve from the caller's op directory. Git refs keep the existing
`repo#tag/path` spelling. `Inputs` snapshots and replaces the supplied arguments;
omitted arguments use the callee's defaults. Values can be handles or JSON
literals, including nested arrays and string-keyed maps. Literal dollar signs
are escaped, and an explicit empty string stays empty.

`PullCreds(username, password)` accepts literal strings or `String` handles.
Without it, the node uses stored auths. Secret interop is temporarily rejected
before execution while private output collection is being added.

Use `call, value := OpOutput[opctl.Dir](call, "artifact")` to request a named
output. Keep the returned call in the graph. Other supported handles include
`String`, `Number`, `Boolean`, `File`, `Socket`, `Array[T]`, and `Object[T]`.
Repeated requests for the same name and type reuse the binding. The SDK checks
producer visibility, including nested input handles and pull credentials.
Result values cannot be read during plan. The callee validates its outputs at
execution; the accessor does not describe or fetch the callee, and existing
opspec coercions still apply.
