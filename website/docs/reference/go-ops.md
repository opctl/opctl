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

Every input field must be exported, use a handle type, and have a unique `opctl` tag name containing letters, numbers, `-`, or `_`. Tags accept `description`, `default`, and `secret`. Boolean inputs cannot be secret; sockets cannot have defaults. Single-quoted tag values can contain commas and escape single quotes or backslashes.

String constraints include `enum`, `format`, `pattern`, `minLength`, and `maxLength`. Number constraints include `enum`, `format=integer`, `minimum`, `maximum`, and `multipleOf`. Arrays accept `minItems`, `maxItems`, and `uniqueItems`. Objects accept `minProperties`, `maxProperties`, and `required=a|b`. Use `constraints='<JSON object>'` for nested JSON Schema constraints. Duplicate keys and unknown tag keys are errors.

Use `EnvFrom` to bind handles, `Dir` and `File` to mount inputs, and `Socket` to bind a socket. Secret input values cannot be read during plan. Directory, file, and socket handles have no `Value()` method. Bootstrap receives only non-secret scalar and structured values; filesystem inputs remain handles. Literal strings escape `$` so input text cannot introduce scope references.

`Describe[Inputs]()` returns the interface without invoking the plan function. `Plan(fn, inputs)` builds the graph in process with already resolved public input values. No-input ops use `func(struct{}) opctl.Call`. Each graph currently contains one container call.

Run a Go op with the same arguments as a YAML op:

```sh
opctl run -a env=prod -a token=example path/to/op
opctl op validate path/to/op
opctl ls path/to/op
```

The CLI prompts for missing inputs, applies defaults, and checks constraints. Invalid explicit Go input values prompt again or fail in a non-interactive terminal; they do not fall back to defaults. Secret prompts are masked. Directory and file defaults such as `$(./testdata)` resolve relative to the op directory. `opctl op validate` checks the declared interface without invoking the plan function. `opctl ls` includes Go ops and their descriptions. Opening `op.go` in the webapp shows its described name and description; its run graph is built when the op runs.

The node bootstraps describe and plan in containers using the Go version in the nearest `go.mod`. Local `replace` targets share a mount that preserves host paths, including absolute replacements. Bootstrap calls appear beneath the op call. Bootstrap stdout is captured in memory, with Docker logging disabled, and is never published as log events. Bootstrap failures return up to 64 KiB of stderr diagnostics to the client. Describe is also available through `POST /api/ops/describes` with `{"op":{"ref":"path/to/op"}}`.

Describe and plan compile separately for each request. No compiled-program cache is available yet. Go ops currently support one container call and no outputs.
