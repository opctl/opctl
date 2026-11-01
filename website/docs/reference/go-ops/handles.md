---
title: Handles
---

A handle is a typed value in the op graph that an input supplies or a call produces: `Dir`, `File`, `Socket`, `String`, `Number`, `Boolean`, `Array[T]`, or `Object[T]`.
Handles that a call produces have no value during plan; [public inputs](inputs-and-outputs.md#values-during-plan) and literal handles do.

## Literal handles

`StringOf`, `NumberOf`, `BooleanOf`, `ArrayOf`, and `ObjectOf` create literal handles. Their `Value()` methods work during plan. Arrays and objects snapshot
their values as JSON. Structs use their `json` tags. Pass these handles to op inputs or `EnvFrom`; strings containing `$()` remain literal text.
Use output default tags for constant op outputs. Literal handles cannot bind an op's exported outputs.

```go
func plan(struct{}) opctl.Call {
    return opctl.CallOp("../greet").Inputs(map[string]any{"name": opctl.StringOf("Ada")})
}
```

## Visibility

A call can use a handle only after its producer has run. The SDK checks this when it builds the op graph, the same way opctl scopes YAML ops:

- In `Serial`, later calls see earlier calls' outputs.
- In `Parallel`, siblings can't see each other's outputs. Their outputs become visible to calls after the group.
- A handle has one producer. Two parallel siblings can't produce the same handle.

Keep every call that a handle accessor returns in the graph: that call produces the handle.

## New directories and files

`NewDir()` and `NewFile()` create handles that a container produces by mounting them. The first container that mounts the handle produces it, and later calls mount the same
directory or file:

```go
func plan(struct{}) opctl.Call {
    out := opctl.NewDir()
    build := opctl.Container("alpine:3").Cmd("sh", "-ce", "echo hello > /out/result.txt").Dir("/out", out)
    check := opctl.Container("alpine:3").Cmd("cat", "/in/result.txt").Dir("/in", out)
    return opctl.Serial(build, check)
}
```

- Create these handles inside the plan function, including in helpers and loops. Their generated names are stable when the graph is built in the same order.
- Mounts are shared and mutable. Parallel calls can mount a handle that is already visible, but their writes need coordination.
- `EnvFrom` can't read a new file before an earlier call mounts it.
- Prefer directory handles to file handles. A tool that replaces a file by renaming another over it can't replace a mounted file. Writing the file in place, for example
  with shell `>`, works.

## Paths in the op's source

`OpDir(path)` and `OpFile(path)` are handles to paths in the op's source, relative to the op directory. The path must start with `./` or `../`.

`Dir(path)` and `File(path)` on a directory handle return handles beneath it. They work on input directories and op directories:

```go
type Inputs struct {
    Home opctl.Dir `opctl:"HOME"`
}

func plan(in Inputs) opctl.Call {
    return opctl.Container("golang:1.25").
        Dir("/src", opctl.OpDir("../..")).
        Dir("/go/pkg", in.Home.Dir("go/pkg")).
        File("/etc/app.yml", opctl.OpDir("./config").File("app.yml")).
        Cmd("true")
}
```

Subpaths are relative, and can't contain `..` elements, `$`, `(`, or `)`, so input text used as a subpath can't leave the directory or add references. Subpaths of call
outputs, such as `NewDir` or `OpOutput` handles, aren't supported yet.

## Op outputs

`OpOutput[T](call, name)` binds an output of an op call to a handle of type `T`, and returns the updated call:

```go
func plan(struct{}) opctl.Call {
    build, artifact := opctl.OpOutput[opctl.Dir](opctl.CallOp("../build"), "artifact")
    return opctl.Serial(build, opctl.Container("alpine:3").Dir("/artifact", artifact).Cmd("ls", "/artifact"))
}
```

Requesting the same name with the same type again returns the same handle; a different type is an error. `T` states how this graph uses the value. The callee checks the
value against its own declaration when it runs, and opctl's usual type coercions apply.

## Stdout and stderr

`Stdout()` and `Stderr()` return the updated container call and a String handle for the captured stream:

```go
func plan(struct{}) opctl.Call {
    producer, version := opctl.Container("alpine:3").
        Cmd("sh", "-ce", `printf 'v1.2.3\n'`).Stdout()
    return opctl.Serial(producer,
        opctl.Container("alpine:3").EnvFrom("VERSION", version).
            Cmd("sh", "-ce", `test "$VERSION" = v1.2.3`))
}
```

Calling `Stdout()` again on the updated call returns the same handle. Captures follow the [container stream
rules](../opspec/op-directory/op/call/container/index.md#stdout-and-stderr), including the 1 MiB limit. Stream handles need opspec schema 2 on the node.

## Outputs file

A container writes named values as a JSON object to the [outputs file](../opspec/op-directory/op/call/container/outputs.md) at `$OPCTL_OUTPUTS`. `ContainerOutput[T](call,
name)` binds one member to a handle:

```go
func plan(struct{}) opctl.Call {
    producer := opctl.Container("alpine:3").Cmd("sh", "-ce",
        `printf '%s' '{"version":"v1.2.3","count":3}' > "$OPCTL_OUTPUTS"`)
    producer, version := opctl.ContainerOutput[opctl.String](producer, "version")
    producer, count := opctl.ContainerOutput[opctl.Number](producer, "count")
    return opctl.Serial(producer, opctl.Container("alpine:3").
        EnvFrom("VERSION", version).EnvFrom("COUNT", count).
        Cmd("sh", "-ce", `echo "$VERSION $COUNT"`))
}
```

`T` is `String`, `Number`, `Boolean`, `Array[T]`, or `Object[T]`. Requesting a member again on the updated call returns the same handle; a different type is an error. `T`
doesn't declare a runtime type: opctl checks only that the member exists and is valid JSON. Use `Outputs` to check types. Outputs file members need opspec schema 3 on the
node.

## Typed container outputs

`Outputs[T](call)` declares outputs file members with a tagged struct of handles, and returns the updated call and the struct:

```go
type BuildResults struct {
    Count    opctl.Number        `opctl:"count,minimum=1,format=integer"`
    Ready    opctl.Boolean       `opctl:"ready"`
    Tags     opctl.Array[string] `opctl:"tags,minItems=1"`
    Metadata opctl.Object[struct {
        Revision string `json:"revision"`
    }] `opctl:"metadata,required=revision"`
}

func plan(struct{}) opctl.Call {
    producer, out := opctl.Outputs[BuildResults](
        opctl.Container("alpine:3").Cmd("sh", "-ce",
            `printf '%s' '{"count":3,"ready":true,"tags":["linux"],"metadata":{"revision":"abc123"}}' > "$OPCTL_OUTPUTS"`),
    )
    consumer := opctl.Container("alpine:3").EnvFrom("COUNT", out.Count).
        Cmd("sh", "-ce", `test "$COUNT" = 3`)
    return opctl.Serial(producer, consumer)
}
```

- `T` is a nonempty struct of exported `String`, `Number`, `Boolean`, `Array`, or `Object` handles. Tags use the [input tag grammar](inputs-and-outputs.md#tags), without
  defaults. Use mounts and socket bindings for filesystem values.
- opctl checks every member's JSON type and constraints before it publishes any of them. A missing member, or a number written as a JSON string, fails the call, and later
  calls don't run.
- `secret` keeps a member's value out of public events. The container must still avoid printing it.
- Declaring the same member again with the same type is allowed; a conflicting declaration is an error. `ContainerOutput` and `Outputs` can declare members on the same
  call.
- Typed container outputs need opspec schema 4 on the node.

## RawRef

`RawRef[T](expr)` emits a whole opspec `$(...)` reference verbatim, for references the SDK doesn't model, such as nested references:

```go
type Inputs struct {
    Configs opctl.Dir    `opctl:"configs"`
    Env     opctl.String `opctl:"env"`
}

func plan(in Inputs) opctl.Call {
    return opctl.Container("alpine:3").
        File("/config.yml", opctl.RawRef[opctl.File]("$(configs/$(env).yml)")).
        EnvFrom("ENV", in.Env).
        Cmd("cat", "/config.yml")
}
```

`RawRef` skips the visibility checks, and its value can't be read during plan. Building the expression from input text can inject references into the graph, so prefer typed
handles.
