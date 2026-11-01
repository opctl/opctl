---
title: Authoring ops in Go
sidebar_label: Authoring ops in Go
---

A Go op is an op declared by a Go program instead of an op file. The program builds an op graph with the Go op SDK, and opctl runs that graph the same way it runs a YAML op. Go ops and YAML ops call each other.

This guide builds up a Go op one feature at a time. The [Go op reference](../reference/go-ops/index.md) has the exact rules for each feature.

## Create and run a Go op

Create a Go op named `hello` in `.opspec`:

```sh
opctl op create --go -d 'Greets the caller.' hello
```

This writes `.opspec/hello/op.go` and `.opspec/hello/go.mod`:

```go title=".opspec/hello/op.go"
// Greets the caller.
package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

type Inputs struct {
	Who opctl.String `opctl:"who,description='Who to greet',default=world"`
}

func plan(in Inputs) opctl.Call {
	return opctl.Container("alpine:3").
		Cmd("sh", "-ce", `echo "hello $WHO"`).
		EnvFrom("WHO", in.Who)
}

func main() { opctl.Main(plan) }
```

If `opctl op create` reports that this build has no published Go op SDK, point `go.mod` at a checkout of opctl before you run the op:

```text title=".opspec/hello/go.mod"
replace github.com/opctl/opctl/sdks/go/opspec/sdk => /path/to/opctl/sdks/go/opspec/sdk
```

Run it like any other op:

```sh
opctl run hello              # hello world
opctl run -a who=opctl hello # hello opctl
```

`op.go` makes the directory a Go op, and the nearest `go.mod` selects the Go version. To run the op, opctl first bootstraps the program in a `golang` container. Bootstrap appears in the event stream and the UI as calls beneath the op.

## Name and description

The op's name is its directory name, and its description is the package doc comment. `opctl ls` and the UI show both. To override either, pass an `opctl.Op` to `Main`:

```go
// Deploys the API.
//
// The description is markdown, so it can hold paragraphs and lists.
package main

import opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"

func plan(struct{}) opctl.Call {
	return opctl.Container("alpine:3").Cmd("echo", "deploying")
}

func main() {
	opctl.Main(plan, opctl.Op{Name: "deploy-api"})
}
```

Go ops have no version field. Version them with git tags, and callers pick a tag in their op ref.

## Inputs

Each field of the inputs struct declares one input. The `opctl` tag holds its name, description, default, and constraints:

```go
type Inputs struct {
	Env      opctl.String `opctl:"env,description='Target environment',default=dev,enum=dev|prod"`
	Replicas opctl.Number `opctl:"replicas,default=1,minimum=1,format=integer"`
	Token    opctl.String `opctl:"token,secret"`
}

func plan(in Inputs) opctl.Call {
	deploy := opctl.Container("alpine:3").
		Cmd("sh", "-ce", `echo "deploying $REPLICAS replicas"`).
		EnvFrom("REPLICAS", in.Replicas).
		EnvFrom("TOKEN", in.Token)
	if in.Env.Value() == "prod" {
		deploy = deploy.Env("STRICT", "true")
	}
	return deploy
}
```

`opctl run` prompts for missing inputs and checks constraints before the op runs. `Value()` returns a public input's value, so Go control flow can branch on it. A secret's value can't be read during plan; bind it to a container with `EnvFrom`.

## Composing calls

`Serial` runs calls in order, and `Parallel` runs them at the same time. Build groups with ordinary Go loops:

```go
func plan(struct{}) opctl.Call {
	var checks []opctl.Call
	for _, check := range []string{"lint", "test", "vet"} {
		checks = append(checks, opctl.Container("golang:1.25").Cmd("echo", check))
	}
	deploy := opctl.Container("alpine:3").Cmd("echo", "deploy")
	return opctl.Serial(opctl.Parallel(checks...), deploy)
}
```

Building the op graph runs nothing. opctl runs the calls after `plan` returns. Build the graph the same way every time for the same inputs: don't depend on map order, time, or goroutines.

## Handles

A handle stands for a value that an input supplies or a call produces. `NewDir` creates a directory handle. The first container that mounts it produces it, and later calls mount the result:

```go
type Inputs struct {
	Home opctl.Dir `opctl:"HOME,description='Home directory of the caller'"`
}

func plan(in Inputs) opctl.Call {
	out := opctl.NewDir()
	build := opctl.Container("golang:1.25").
		Dir("/src", opctl.OpDir("../..")).
		Dir("/go/pkg", in.Home.Dir("go/pkg")).
		Dir("/out", out).
		WorkDir("/src").
		Env("CGO_ENABLED", "0").
		Cmd("go", "build", "-o", "/out/app", ".")
	test := opctl.Container("alpine:3").
		Dir("/in", out).
		File("/check.sh", opctl.OpFile("./check.sh")).
		Cmd("sh", "/check.sh", "/in/app")
	return opctl.Serial(build, test)
}
```

`OpDir` and `OpFile` are handles to paths in the op's source, relative to the op directory. `in.Home.Dir("go/pkg")` is a handle to a directory beneath an input. As with any input, `opctl run` takes `HOME` from an environment variable of the same name. When no handle fits, `RawRef` emits an opspec `$(...)` reference verbatim, without any checks.

## Helpers

A helper is a plain Go function that adds calls to the caller's op graph. It returns its calls and any handles they produce:

```go
// goBuild builds the module in src and returns the directory it fills.
func goBuild(src opctl.Dir) (opctl.Call, opctl.Dir) {
	out := opctl.NewDir()
	return opctl.Container("golang:1.25").
		Dir("/src", src).
		Dir("/out", out).
		WorkDir("/src").
		Cmd("go", "build", "-o", "/out/", "./..."), out
}

func plan(struct{}) opctl.Call {
	build, bin := goBuild(opctl.OpDir("../.."))
	return opctl.Serial(build, opctl.Container("alpine:3").Dir("/app", bin).Cmd("ls", "/app"))
}
```

A helper has no op boundary: its calls are part of your op. To reuse a whole op, with its own inputs and outputs, call it by ref instead.

## Calling YAML ops, and being called

`CallOp` calls another op by ref, whether it's a YAML op or a Go op. `OpOutput` binds one of its outputs to a handle:

```go
func plan(struct{}) opctl.Call {
	version, tag := opctl.OpOutput[opctl.String](
		opctl.CallOp("../version").Inputs(map[string]any{"prefix": "v"}),
		"tag",
	)
	publish := opctl.Container("alpine:3").
		EnvFrom("TAG", tag).
		Cmd("sh", "-ce", `echo "publishing $TAG"`)
	return opctl.Serial(version, publish)
}
```

Local refs are relative to the op directory. Git refs use `repo#tag/path`. A YAML op calls a Go op the same way it calls any op:

```yaml
name: release
run:
  op:
    ref: $(../build)
    inputs:
      env: prod
```

## Returning op outputs

Declare outputs with a second struct, and use `MainWithOutputs`. The plan function returns the root call and the handles to export:

```go
type Outputs struct {
	Binary  opctl.Dir    `opctl:"binary,description='Built executables'"`
	Channel opctl.String `opctl:"channel,default=stable"`
}

func plan(struct{}) (opctl.Call, Outputs) {
	out := opctl.NewDir()
	build := opctl.Container("golang:1.25").
		Dir("/src", opctl.OpDir("../..")).
		Dir("/out", out).
		WorkDir("/src").
		Cmd("go", "build", "-o", "/out/", "./...")
	return build, Outputs{Binary: out}
}

func main() { opctl.MainWithOutputs(plan) }
```

Every output needs a bound handle or a default. Here `Channel` keeps its default.

## Stdout and the outputs file

A container can return values without writing files to a mount. `Stdout` captures its standard output as a String handle. A container can also write JSON to the outputs file at `$OPCTL_OUTPUTS`, and `ContainerOutput` binds one of its members:

```go
func plan(struct{}) opctl.Call {
	read, version := opctl.Container("alpine:3").
		File("/VERSION", opctl.OpFile("./VERSION")).
		Cmd("cat", "/VERSION").
		Stdout()
	counter, count := opctl.ContainerOutput[opctl.Number](
		opctl.Container("alpine:3").Cmd("sh", "-ce", `echo '{"count": 3}' > "$OPCTL_OUTPUTS"`),
		"count",
	)
	report := opctl.Container("alpine:3").
		EnvFrom("VERSION", version).
		EnvFrom("COUNT", count).
		Cmd("sh", "-ce", `echo "$VERSION has $COUNT checks"`)
	return opctl.Serial(opctl.Parallel(read, counter), report)
}
```

Keep the returned call in the graph: it produces the handle. Calls that use a handle must run after its producer, so `report` runs after the parallel group.

## Typed container outputs

`Outputs` declares the members of the outputs file with a tagged struct. opctl checks each member's type and constraints when the container exits:

```go
type BuildResults struct {
	Count opctl.Number        `opctl:"count,minimum=1,format=integer"`
	Tags  opctl.Array[string] `opctl:"tags,minItems=1"`
}

func plan(struct{}) opctl.Call {
	build, out := opctl.Outputs[BuildResults](
		opctl.Container("alpine:3").Cmd("sh", "-ce",
			`echo '{"count": 3, "tags": ["linux"]}' > "$OPCTL_OUTPUTS"`),
	)
	check := opctl.Container("alpine:3").
		EnvFrom("COUNT", out.Count).
		Cmd("sh", "-ce", `test "$COUNT" -ge 1`)
	return opctl.Serial(build, check)
}
```

A missing member or a value of the wrong type fails the call, so `check` doesn't run.

## Modules and app packages

A Go op doesn't need its own module. Delete the generated `go.mod`, and require the Go op SDK in your repository's `go.mod`. The op can then import your app's packages:

```text
myapp/
├── go.mod                    module example.com/myapp
├── internal/deploy/deploy.go
└── .opspec/deploy/op.go
```

```go title="internal/deploy/deploy.go"
// Package deploy holds settings shared by the app and its ops.
package deploy

// Image is the container image that runs the API.
const Image = "alpine:3"
```

```go title=".opspec/deploy/op.go"
// Deploys the API.
package main

import (
	"example.com/myapp/internal/deploy"
	opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"
)

func main() {
	opctl.Main(func(struct{}) opctl.Call {
		return opctl.Container(deploy.Image).Cmd("echo", "deploying")
	})
}
```

Bootstrap mounts the module read-only, so update `go.mod` and `go.sum` outside opctl. An op fetched by git ref must keep its `go.mod` and every local `replace` target inside the fetched repository.

## Caching

opctl compiles a Go op once and reuses the program:

```sh
opctl run deploy   # compiles the program in a bootstrap build call
opctl run deploy   # reuses it: no bootstrap build call
```

The program recompiles when the module's source, `go.mod`, `go.sum`, local replacements, or the Go image change. Plan runs again for every call, and input values are never cached.

## Private modules

Bootstrap downloads private modules with the auths you already added to opctl:

```sh
opctl auth add github.com/acme -u x-access-token -p "$GITHUB_TOKEN"
```

With that auth, an op can require `github.com/acme/tools` in its `go.mod`. opctl derives `GOPRIVATE` from your auths, and only module downloads see the credentials.

## Testing

`Plan` builds the op graph in a Go test, without a node. Compare it with a saved copy:

```go title=".opspec/hello/op_test.go"
package main

import (
	"bytes"
	"os"
	"testing"

	opctl "github.com/opctl/opctl/sdks/go/opspec/sdk"
)

func TestPlan(t *testing.T) {
	got, err := opctl.Plan(plan, map[string]any{"who": "opctl"})
	if err != nil {
		t.Fatal(err)
	}
	want, err := os.ReadFile("testdata/plan.json")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("plan = %s, want %s", got, want)
	}
}

func TestDescribe(t *testing.T) {
	if _, err := opctl.Describe[Inputs](); err != nil {
		t.Fatal(err)
	}
}
```

`Plan` takes resolved values, so pass every public input, including those with defaults. `Describe` returns the op's interface without calling `plan`.

## Next

- The [Go op reference](../reference/go-ops/index.md) covers every SDK function and the rules opctl applies.
- The [run a go service](../run-a-go-service.md) example has a Go op that calls a YAML op.
