---
title: Inputs and outputs
---

## Inputs

The inputs struct declares an op's inputs. Each exported field is one input, typed with a handle, and its `opctl` tag names it:

```go
type Inputs struct {
	Env    opctl.String  `opctl:"env,description='Target environment',default=dev,enum=dev|prod"`
	Count  opctl.Number  `opctl:"count,default=2,minimum=1,maximum=3,format=integer"`
	Token  opctl.String  `opctl:"token,secret,minLength=1"`
	Src    opctl.Dir     `opctl:"src,default=$(./testdata/src)"`
	DryRun opctl.Boolean `opctl:"dryRun,default=false"`
}
```

| Handle | Input type | `Value()` during plan |
| --- | --- | --- |
| `String` | string | `string` |
| `Number` | number | `float64` |
| `Boolean` | boolean | `bool` |
| `Array[T]` | array | `[]T` |
| `Object[T]` | object | `T`, a struct or string-keyed map |
| `Dir` | dir | none |
| `File` | file | none |
| `Socket` | socket | none |

`Array` and `Object` derive their item and property schemas from `T`, using `json` tags for property names. `[]any` and `map[string]any` allow untyped values. Unsigned integer types add `minimum: 0`. Recursive types, pointers, and embedded struct fields are rejected.

Every field must be exported, use a handle type, and have a unique name of letters, numbers, `-`, and `_`. Names starting with `opctl_` are reserved.

### Tags

A tag starts with the input's name, followed by comma-separated keys:

| Key | Applies to | Meaning |
| --- | --- | --- |
| `description` | all | Markdown description |
| `default` | all but `Socket` | Default value; dir and file defaults such as `$(./testdata)` are relative to the op directory |
| `secret` | all but `Boolean` | The value stays out of plan and of public events |
| `enum` | `String`, `Number` | Allowed values, separated by `\|` |
| `format`, `pattern`, `minLength`, `maxLength` | `String` | String constraints |
| `format=integer`, `minimum`, `maximum`, `multipleOf` | `Number` | Number constraints |
| `minItems`, `maxItems`, `uniqueItems` | `Array` | Array constraints |
| `minProperties`, `maxProperties`, `required=a\|b` | `Object` | Object constraints |
| `constraints='<JSON object>'` | all | Extra JSON Schema keywords |

Quote a value with single quotes to include commas, and escape a single quote or backslash with a backslash. Duplicate and unknown keys are errors. `constraints` can add keywords, but can't replace those derived from the Go type.

### Values during plan

`Value()` returns a public input's resolved value, so Go control flow can branch on it. Reading a secret input fails the plan; bind it to a container with `EnvFrom` instead. Directory, file, and socket handles have no value during plan.

Bootstrap receives only public string, number, boolean, array, and object values. Secret and filesystem inputs reach the op graph as handles.

## Outputs

Use `MainWithOutputs` to return outputs. The plan function returns the root call and a struct of output handles. Output tags use the same grammar as input tags:

```go
type Inputs struct {
	Text opctl.String `opctl:"text"`
}

type Outputs struct {
	Echoed   opctl.String `opctl:"echoed"`
	Fallback opctl.String `opctl:"fallback,default=from-default"`
}

func main() {
	opctl.MainWithOutputs(func(in Inputs) (opctl.Call, Outputs) {
		return opctl.Container("alpine:3").Cmd("true"), Outputs{Echoed: in.Text}
	})
}
```

- Every output must hold a bound handle or have a default. A zero handle uses the default. An invalid handle, or one no call produces, fails even with a default.
- An output can bind an input, a [call output](handles.md), or a container result. Two outputs can bind the same handle, and an input and an output can share a name.
- Describe reads only the tagged types, and never calls the plan function.
- opctl checks the bindings against the described outputs before it runs the op graph, then resolves them when the graph completes. A failed binding names the output and the underlying error; for a secret binding, it names only the output.

API clients receive a root Go op's public outputs in its completion event. Nested calls export the outputs their caller requests. Secret outputs never appear in public events. Dir and file outputs are paths on the node.
