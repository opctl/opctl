---
title: Op stubs
---

A typed op call uses an op stub generated from a callee's describe interface. The callee can be a local or git-ref YAML or Go op. Ops stay in `.opspec/`.
The stub is an ordinary package in the caller's Go module, at `opctlgen/NAME`, importing only the Go op SDK.

## Declare an op stub

Add a trailing `//opctl:op REF` comment to its import. For a repository module `example.com/myapp`:

```go title=".opspec/typed-example/op.go"
package main

import (
    "example.com/myapp/opctlgen/greet" //opctl:op ../greet
    opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"
)

func main() {
    opctl.Main(func(struct{}) opctl.Call {
        call, out := greet.Call(greet.Inputs{
            Name: opctl.StringOf("Ada"),
            Tags: opctl.ArrayOf([]string{"release"}),
            Config: opctl.ObjectOf(greet.InputsConfig{Style: "loud", Loud: true}),
        })
        return opctl.Serial(call, opctl.Container("alpine:3").EnvFrom("GREETING", out.Greeting).Cmd("true"))
    })
}
```

- Exactly one ref follows the comment. Write `//opctl:op`, without a space after `//`, on the import itself.
- Local refs start with `./` or `../` and are relative to the declaring file. A fetched op's local target must stay inside its fetched repository.
- Git refs use the existing `repo#tag/path` grammar. The tag pins the ref. A refreshed tag's changed interface invalidates the caller's compiled program.
- Names match `[a-z][a-z0-9]*`. Go keywords, `main`, `internal`, `vendor`, and `testdata` are reserved. Stub paths are one level deep.
- Declarations are module-wide. A bare import is allowed when another file declares the same stub. Repeated declarations must resolve to the same target.
- Test files count for `op gen`, but their imports do not add stubs to an op's build. Bootstrap follows helper imports and scans across build constraints.

## Inputs and outputs

Each package has `Inputs`, `Outputs`, and `Call(Inputs) (opctl.OpCall, Outputs)`. Keep the returned call in the graph. Its outputs follow the ordinary
[visibility rules](handles.md#visibility). `.PullCreds(...)` can still be chained on the returned call.

| Callee parameter | Generated field |
| --- | --- |
| string, number, boolean | `opctl.String`, `opctl.Number`, `opctl.Boolean` |
| dir, file, socket | `opctl.Dir`, `opctl.File`, `opctl.Socket` |
| array | `opctl.Array[T]`; unknown items use `any` |
| object with properties | `opctl.Object[InputsField]` or `opctl.Object[OutputsField]`, with generated nested structs |
| object with an additional-properties schema | `opctl.Object[map[string]T]` |
| object without a schema | `opctl.Object[map[string]any]` |

Property schemas map strings to `string`, integers to `int`, numbers to `float64`, booleans to `bool`, and arrays and objects recursively. Unsupported
subtrees, including unions, references, and null, use `any`. The callee still validates values and constraints when it runs.

Leave a handle unset to use the callee's default. An input without a default is required; an unset required input fails during plan. Secret defaults never
appear in generated files. Descriptions and constraints become comments.

Optional object properties use `omitempty`. An explicit zero value for an optional property cannot be sent through its generated struct. Use `CallOp` with
a map when you need to distinguish zero from omission. Typed object handles with different generic arguments need an explicit Go conversion when compatible.

## Names and literal handles

Fields split parameter names on non-ASCII-alphanumeric characters and capitalize each part. `db-name` becomes `DbName`; `api_URL` becomes `ApiURL`.
Leading digits get `X`. Names without ASCII letters or digits get `X` plus eight SHA-256 hex digits. Every colliding name gets an eight-digit hash suffix,
so adding a collision causes a compiler error instead of silently moving a caller's input to another parameter. Nested type names use the same collision rule.

`StringOf`, `NumberOf`, `BooleanOf`, `ArrayOf`, and `ObjectOf` make literal handles. Objects and arrays snapshot JSON, respecting struct `json` tags.
`Value()` reads these literals during plan. Strings containing `$()` stay literal. Literal handles work in op inputs and `EnvFrom`; `StringOf` also works
in pull credentials. There are no literal directory, file, or socket constructors. Use `OpDir` and `OpFile` for source paths. Export constants with output
default tags rather than literal output bindings.

## Generate, run, and edit

`opctl op gen [OP_REF...]` refreshes every declaration in the local refs' modules. With no arguments it finds Go ops under `./.opspec`. It reports every
generation failure and exits nonzero, while writing successful stubs. Local `run` and `op validate` refresh the whole module and warn on failure.

`opctlgen/.gitignore` contains `*`, and Go op scaffolds also ignore `/opctlgen/`. No generated file is committed. Refresh refuses user files and symlinks.
Missing stubs cause editor import errors; generation fixes them. Helpers outside `.opspec/` may also need generation before `go mod tidy`.

The node always renders execution stubs from current callee interfaces and builds through a read-only Go overlay. It hides host stub files and excludes
them from source fingerprints. A changed callee interface forces a new compiled program. Plan output and input values are never cached.

## Diagnostics

Errors name the stub and its declaring file and line. Fix the declaration or caller, then run `opctl op gen` for editor completion.

| Diagnostic | Action |
| --- | --- |
| `no ref is declared for op stub` | Add the trailing import comment. |
| `needs one ref` or `without a space` | Use exactly `//opctl:op REF`. |
| `lowercase Go package name` or `one level` | Use an allowed name directly below `opctlgen/`. |
| `use ./ or ../` | Supply a relative local ref or `repo#tag/path`. |
| `declared with different refs` | Use one target per name, or rename one import. |
| `no op at` | Check the ref relative to its declaring file. |
| `outside fetched repository root` | Keep local callees in the fetched repository. |
| `fetch` or authentication failure | Check the git ref and add credentials with `opctl auth add`. |
| `describe` with an op declaration error | Fix the callee's interface. |
| `op stubs form a cycle` | Use `CallOp` for one edge of the Go-to-Go cycle. YAML interfaces end recursion. |
| `unknown field`, missing output, or changed type | Update the caller to the current interface. The compiler hint lists stubs and marks stale editor files. |
| `undefined: opctl.StubCall` | Update the SDK requirement to the release matching opctl. |
| `input "name" is required` | Set that input handle or add a default to the callee. |

## Version compatibility

Typed op calls require the opctl release containing this feature, currently **Unreleased**, and its matching Go op SDK. Earlier opctl releases cannot
generate these packages. They may fail with a missing `opctlgen/` import when bootstrapping fetched source. Stubs add no contract fields or opspec schema
version; their emitted calls use the existing schema 1 op-call representation.
