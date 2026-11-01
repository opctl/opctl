---
title: Bootstrap
---

Op stub imports carry a trailing `//opctl:op <ref>` comment. The node describes each callee and renders the stub before compiling the caller. A read-only
overlay supplies those packages to Go. Changes to a callee's interface change the caller's compiled-program cache key. Host `opctlgen/` files never supply
execution code and do not affect the cache key. Go-to-Go stub cycles fail with a suggestion to use `CallOp` for one edge.

Bootstrap is how opctl obtains a Go op's declaration: it compiles the program and runs its [describe and plan phases](index.md#describe-and-plan) in containers
it manages. Bootstrap appears in the event stream and the UI as calls beneath the op call: `bootstrap build`, `bootstrap describe`, and `bootstrap plan`.

## Go version and image

The `go` directive in the nearest `go.mod`, found by walking up from the op directory, selects the `golang:<version>` image. opctl resolves the image to an
immutable local image before bootstrap, so a cached versioned image works offline. Bootstrap runs Linux containers on every host.

## Source mounts

- The module root and each local `replace` target outside it mount read-only at their host paths, including absolute replacements. Replacements inside the
  module use the module's mount.
- Bootstrap builds the whole op package, so an op in a subdirectory can import its module's packages and embedded files.
- Mounts can't follow host symlinks, and can't hide the image's `/usr`, `/usr/local`, `/go`, or `/opctl` directories. Missing replacement targets aren't
  mounted; Go reports them if the build needs them.
- Bootstrap can't write its source mounts. Compilation uses `-mod=readonly -buildvcs=false -trimpath`, so update `go.mod` and `go.sum` outside opctl when Go
  reports missing dependencies.
- An op fetched by git ref must keep its `go.mod` and every local `replace` target inside the fetched repository. Bootstrap rejects paths and symlinks that
  leave it. Local ops can use any parent module and replacement.
- With Docker, bootstrap needs Docker API version 1.44 or later, which enforces recursive read-only mounts.

## Build environment

Bootstrap runs Go as the owner of the node's Go caches, with a temporary home and `GOWORK=off`, `GOENV=off`, `GOTOOLCHAIN=local`, and `CGO_ENABLED=0`. A parent
`go.work` and your own Go environment don't affect it.

## Caches

| Cache | Location | Contents |
| --- | --- | --- |
| Modules | `<dataDir>/go/mod` | Downloaded module source, shared by all ops and Go versions |
| Builds | `<dataDir>/go/build` | Go's build cache |
| Programs | `<dataDir>/go/programs/v1` | Compiled Go op executables |
| Describe | node memory | Successful describe responses |

A compiled program is keyed by the module's source content, local replacements, op package, `go` directive, image, and build settings. Uncommitted and untracked
files, embedded files, and edits to a replaced SDK change the key; changed modification times alone don't. The program cache has a 1 GiB soft limit and evicts
the least recently used programs that no call is using. Concurrent calls share one build, and a failed or cancelled build stores nothing.

Describe responses are kept in node memory, up to 64 MiB or 128 entries. A cached describe appears as a `bootstrap describe (cached)` call that starts no
container. Restarting the node clears it. Plan runs for every call, and neither plan responses nor input values are ever stored.

To clear a cache, stop the node and delete its directory. Deleting `<dataDir>/go` removes the modules, builds, and programs. With the embedded container
runtime, the data directory must be shared with its VM at the same path. The caches aren't yet verified with the embedded runtime on macOS.

## Private modules

Bootstrap downloads private modules with the HTTPS username and password auths added by `opctl auth add`:

```sh
opctl auth add github.com/acme -u x-access-token -p "$GITHUB_TOKEN"
```

- Use a host or repository path, such as `github.com`, `github.com/acme`, or `https://git.example.com/team/repo`. Image refs, wildcards, and SSH resources don't
  configure module auth.
- opctl derives `GOPRIVATE` from these auths and keeps public proxy and checksum verification for other modules. The longest matching path wins, and matches
  respect host and path boundaries: an auth for `github.com/acme/repo` doesn't match `github.com/acme/repository`.
- A fetched Go op's pull credentials override stored auths for modules in its own repository. Other private modules use their own matching auths. A local op's
  pull credentials don't apply to modules; add an auth for its private dependencies. Image pull credentials are separate from module credentials.
- Credentials reach only module downloads, through container memory. They're removed before compilation, and describe and plan never receive them. They're never
  written to the caches, events, or logs.
- Credentials can't contain line breaks. Authenticated vanity import paths that need repository-scoped credentials aren't supported with Go 1.21.

Revoking a credential doesn't remove downloaded source or compiled programs. Delete `<dataDir>/go` to remove them.

## Diagnostics

Bootstrap stdout carries only the protocol response. opctl captures it in memory and never publishes it as log events. When compilation or a declaration fails,
clients receive up to 64 KiB of stderr. Download failures that involve credentials return a fixed error without Go's output.

## Release retries

Each opctl release tags the Go op SDK as `sdks/go/opspec/sdk/vX.Y.Z` before compiling and publishing opctl. A new tag points to the release commit.
If a release fails later, retrying keeps the existing SDK tag when it points to the
release commit or its `sdks/go/opspec/sdk` contents match the release commit. The release fetches both commits to compare them and reuses the tagged module's
`go.sum` hashes. It never moves the tag. If the SDK contents differ, bump the CHANGELOG version before releasing again.

## Describe API

`POST /api/ops/describes` with `{"op":{"ref":"path/to/op"}}` returns a Go or YAML op's declaration. Invalid declarations, including a missing or malformed
`go.mod`, invalid mounts, bootstrap program failures, and unreadable responses, return HTTP 400. Runtime and context failures return HTTP 500.
