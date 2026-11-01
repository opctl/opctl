---
title: Calls
---

Generated op stubs use `StubCall` to resolve local refs relative to the caller's module and check required inputs during plan. Hand-written calls use `CallOp`.

A plan function returns one call, the root of the op graph. Calls are container calls, op calls, and groups. Builders return independent copies, so a call can
be the base of several others.

## Container calls

`Container(image)` starts a container call with a literal image ref.

| Method | Sets |
| --- | --- |
| `Cmd(args...)` | The literal command and arguments |
| `Env(name, value)` | A literal environment variable |
| `EnvFrom(name, handle)` | An environment variable bound to a handle |
| `Dir(path, dir)` | A directory mount |
| `File(path, file)` | A file mount |
| `Socket(address, socket)` | A socket binding |
| `WorkDir(path)` | The literal working directory |
| `Name(name)` | A literal name that other opctl containers and the host resolve through opctl DNS. Containers that share a name are load-balanced. |
| `Stdout()`, `Stderr()` | A String handle for the captured stream; see [Handles](handles.md#stdout-and-stderr) |

Literal strings escape every `$`, so input text can't add opspec references to the graph. Binding names can't contain `$`. Setting a binding twice keeps the last value.

```go
func plan(struct{}) opctl.Call {
    return opctl.Container("golang:1.25").
        Name("api").
        Dir("/src", opctl.OpDir("../..")).
        Env("CGO_ENABLED", "0").
        WorkDir("/src").
        Cmd("go", "run", ".")
}
```

## Op calls

`CallOp(ref)` calls another op, which can be a YAML op or a Go op. A local ref is relative to the op directory. A git ref uses `repo#tag/path`. opctl fetches
the callee when the call runs, not during plan.

```go
type Inputs struct {
    Label opctl.String `opctl:"label"`
}

func plan(in Inputs) opctl.Call {
    return opctl.CallOp("../build").Inputs(map[string]any{
        "label":  in.Label,
        "labels": []any{"a", in.Label},
        "empty":  "",
    })
}
```

- `Inputs` replaces the call's arguments with a snapshot of the map. Values can be handles or JSON-compatible literals, including nested slices and string-keyed
  maps. Omitted arguments use the callee's defaults, and an empty string stays empty.
- `PullCreds(username, password)` sets the credentials for fetching the callee, as literal strings or `String` handles. Without it, opctl uses the auths added
  with `opctl auth add`.
- [`OpOutput`](handles.md#op-outputs) binds one of the callee's outputs to a handle.

Secret values stay out of public events across op calls, including through renamed outputs, groups, and YAML loops. To pass a secret into a Go op's input, that
input must have the `secret` tag.

## Groups

`Serial(calls...)` runs calls in order; later calls see earlier calls' outputs. `Parallel(calls...)` runs calls at the same time; each starts with the group's
incoming handles, and their outputs become visible after the group. Empty groups are dropped from their parent, and an empty root graph is an error.

## Control flow and helpers

Plain Go builds the graph. Use `if` on input values to choose calls, and `for` loops to build a slice of calls for a group:

```go
func checks(names []string) opctl.Call {
    var calls []opctl.Call
    for _, name := range names {
        calls = append(calls, opctl.Container("alpine:3").Cmd("echo", name))
    }
    return opctl.Serial(opctl.Parallel(calls...), opctl.Container("alpine:3").Cmd("echo", "done"))
}
```

A function like `checks` is a helper: it adds calls to the caller's op graph, without an op boundary. Build the graph deterministically in one goroutine, and
sort map keys before iterating over them.
