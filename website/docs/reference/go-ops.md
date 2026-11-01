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

The node bootstraps describe and plan in containers using the Go version in the nearest `go.mod`. It resolves the selected image to an immutable local image identity before running bootstrap. Cached versioned images support offline execution. The module root and each external local `replace` target mount read-only, separately at their host paths, including absolute replacements. Replacements inside the module use the module mount. Separate mounts and module roots cannot follow host symlinks. Missing local replacement targets are not mounted; Go reports an error if the build needs them. Mounts that hide required bootstrap image directories are rejected. Bootstrap cannot write its source mounts or remount them writable. It builds the entire op package, so an op in a subdirectory can import its application module's packages and embedded assets. Docker must support recursive read-only binds. Bootstrap calls appear beneath the op call. Bootstrap stdout is captured in memory, with Docker logging disabled, and is never published as log events. Compiler and declaration failures return up to 64 KiB of stderr diagnostics to the client. Credentialed dependency downloads return controlled errors, with their raw diagnostics suppressed. Describe is also available through `POST /api/ops/describes` with `{"op":{"ref":"path/to/op"}}`. Invalid declarations, including missing or malformed `go.mod` files and invalid mount paths, bootstrap program exits, and response parse failures return HTTP 400; runtime and context failures return HTTP 500. This classification applies to bootstrap; YAML container calls retain their existing exit-code errors.

Describe and plan use the same cached executable. Describe can use the in-memory memo; plan executes for every request.

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
Without it, the node uses stored auths. Secret bindings stay in runtime memory. Public call events omit their values,
including renamed outputs and values carried through groups or YAML loops. A
sensitive value passed into a Go scalar or structured input requires the
`secret` tag on that input.

Use `call, value := OpOutput[opctl.Dir](call, "artifact")` to request a named
output. Keep the returned call in the graph. Other supported handles include
`String`, `Number`, `Boolean`, `File`, `Socket`, `Array[T]`, and `Object[T]`.
Repeated requests for the same name and type reuse the binding. The SDK checks
producer visibility, including nested input handles and pull credentials.
Result values cannot be read during plan. The callee validates its outputs at
execution; the accessor does not describe or fetch the callee, and existing
opspec coercions still apply.

The node validates described Go outputs against runtime bindings before it runs
the graph, then resolves those bindings from the completed graph scope. Root Go
calls export every declared output in their completion event. Nested calls
export the names requested by their caller. Defaults fill omitted bindings;
output declarations never turn missing runtime references into fresh mounts.

## Declaring outputs

Use `MainWithOutputs` when a Go op returns outputs. Its function returns the
root call and a struct of output handles. Inputs and outputs use the same
`opctl` tag grammar, including descriptions, constraints, defaults, and `secret`.

```go
type Inputs struct {
    Text opctl.String `opctl:"text"`
}
type Outputs struct {
    Echoed opctl.String `opctl:"echoed"`
    Fallback opctl.String `opctl:"fallback,default=from-default"`
}
func main() {
    opctl.MainWithOutputs(func(in Inputs) (opctl.Call, Outputs) {
        return opctl.Container("alpine:3").Cmd("true"), Outputs{Echoed: in.Text}
    })
}
```

Every output must hold a bound handle or have a default. A zero handle uses its
default; an invalid or unproduced handle fails even when it has a default. Two
outputs can alias one handle. Input and output declarations may share a public
name. Use `PlanWithOutputs` and `DescribeWithOutputs` in SDK tests. Describe
reads only the tagged types and never invokes the plan function. Existing
`Main`, `Plan`, and `Describe` programs keep their no-output contract.

Go ops fetched by git ref must keep their nearest `go.mod` and every local
`replace` target inside that fetched repository. Bootstrap rejects paths and
symlinks that escape it. The node carries the resolved repository root through
nested local calls; it does not infer the boundary from `#` or a `.git` folder.
Local ops retain their existing parent-module and replacement behavior. A git
fetch includes the tagged repository tree, including helpers and in-repo SDK
replacements, but does not fetch a developer's sibling checkouts.

Public Go outputs reach API clients as the boxed `CallEnded.Outputs` map on the
root op completion event. `opctl run` forwards this event to its output module
and prints the usual status. It does not print a new output format. Secret
outputs and sensitive aliases are omitted from public events. Dir and file
outputs refer to paths on the node, rather than downloaded client files.

Bootstrap shares Go's module cache at `<dataDir>/go/mod` and build cache at
`<dataDir>/go/build` across ops and Go versions. It drops to the mount owner's
UID/GID before invoking Go, with a temporary home and `GOWORK=off`, `GOENV=off`,
`GOTOOLCHAIN=local`, and `CGO_ENABLED=0`. Compilation uses
`-mod=readonly -buildvcs=false -trimpath`; update `go.mod` and `go.sum` outside
bootstrap when Go reports dependency-maintenance errors. An ancestor `go.work`
and ambient build flags do not affect bootstrap.

The module cache retains downloaded source and can grow over time. To remove it,
stop the node and delete `<dataDir>/go/mod`. Caches belong to one node/data-dir
owner. With embedded, the data directory must already have a Lima share at the
same absolute path. Custom unshared paths require configuring that share first.

Compiled programs live in `<dataDir>/go/programs/v1`, keyed by source content,
local replacements, op package, Go directive, immutable image identity, and build
settings. Dirty and untracked files, embedded assets, and SDK replacement edits
invalidate the key; mtime changes alone do not. Escaping source symlinks and
special files fail with a path-specific error. The cache stores only executables
and diagnostic metadata, never describe or plan responses. Every plan executes
again with its own public inputs. Concurrent callers share a build; cancellation
or a failed build never publishes an entry.

Programs use a 1 GiB soft LRU limit; active calls keep their executables leased.
An oversized executable is available only for its current call. To force all
programs to rebuild, stop the node and remove `<dataDir>/go/programs`.

Successful describe responses are memoised only in node memory, with a 64 MiB /
128-entry LRU bound. Each caller receives a freshly parsed declaration. A hit
emits `bootstrap describe (cached)` lifecycle events without starting a
container. Restarting the node retains compiled executables and clears describe
memoisation. Errors and plans are never memoised.

The Docker Linux cache proof, with images prepared and empty node Go caches,
measured 16.014 s cold and 2.971 s warm, a 5.39x improvement. Bootstrap took
14.146 s cold, including 12.304 s building, and 1.317 s warm with no build.
Source hashing took 3.321 ms. The warm run kept the executable checksum and
mtime, skipped describe execution, and planned again with a different input.
These fixture measurements are not a speed guarantee for other modules.

Private Go modules use existing `opctl auth add` HTTPS username/password records.
Use concrete resources such as `github.com`, `github.com/acme`, or
`https://git.example.com/team/repo`. Bootstrap derives `GOPRIVATE` from these
scopes, keeps public proxy/checksum verification enabled during download, and
uses the longest matching host/repository path for Git credentials. Host and path
boundaries prevent an auth for `github.com/acme/repo` from matching
`github.com.evil` or `github.com/acme/repository`. Image tags, versions, wildcards,
SSH resources, and non-host resources do not configure module auth.

Explicit pull credentials for a fetched Go op override stored credentials only
within that source repository's namespace. A local op cannot infer a remote scope
from pull credentials; add an auth record for its private dependencies. Other and
transitive private modules use their own matching records. Image pull credentials
are kept separate from module credentials.

The node prepares credentials in memory and sends a bounded archive through
attached stdin into container tmpfs. Its Git helper preserves literal spaces,
quotes, dollar signs, and backslashes; line breaks and NUL are rejected. Only
host-wide credentials representable by older Go netrc readers are used for HTTP
discovery. Whitespace, quoted, and escaped credentials still work through the Git
helper. Authenticated vanity discovery with repository-only credentials is not
supported on Go 1.21; bootstrap reports that limit without widening the scope.
SSH and newer `GOAUTH` authentication are outside this implementation.

Downloads finish before compilation. Bootstrap then deletes the helper, netrc,
Git config, archive, and download diagnostics from tmpfs, resets its home, and
compiles with module network fallback disabled. Describe and plan receive only
the read-only program/source mounts and their request. No protocol response or
additional credential copy is stored in the node caches, call events, Docker logs,
or Git origin URLs. The existing `AuthAdded` event/auth store already persists the
original credentials; this feature does not change that store.

Revoking or rotating credentials does not erase downloaded source or compiled
code or invalidate an executable. If code retention must end, stop the node and
remove `<dataDir>/go`, including its module, build, and program caches. These
caches are private to one data-dir owner, not an authorization boundary between
users. Docker Linux has integration coverage; embedded requires a same-path Lima
share and still needs a macOS smoke test for UID mapping and read-only behavior.

Container graphs share YAML stdout/stderr bindings and privacy rules. The Go SDK accessors arrive separately; the host reads schema 1 graphs and schema 2 stream bindings.
