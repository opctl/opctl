---
name: converting-opctl-ops
description: Converts opctl ops between op.yml and Go. Use when porting a YAML op to a Go op, finishing a port that `opctl op convert --go` refused, or porting a Go op back to op.yml.
---

# Converting opctl ops

A conversion keeps the op's **interface** exactly: name, description, inputs, outputs, types, defaults, constraints, secrets, and
descriptions. Callers can't tell the op changed; only its declaration changes form.

## YAML to Go

1. **Preview.** From the directory holding `.opspec/`, run `opctl op convert --go <op> --dry-run`. Done when it prints either op.go (plus
   go.mod and .gitignore when the op gets its own module) or a list of **refusals**.
2. **Split refusals.** Skip when nothing was refused. Each refusal names an op.yml location and a construct the Go op SDK can't express.
   Move the smallest call containing it into a new YAML op beside this one, such as `.opspec/<op>-ports/op.yml`: give it an input for every
   `$(...)` the call reads and an output for every name later calls read, and replace the call in the original with an `op:` call to
   `../<op>-ports`. Check the **split** with `opctl op plan` on both ops, then return to step 1. Done when step 1 prints op.go. Some
   refusals have a Go rewrite instead; see [Refusals](#refusals).
3. **Convert.** Run `opctl op convert --go <op>`. It checks that the Go op plans the op.yml graph for the default inputs and generated
   samples, then writes the files, removes op.yml, and refreshes `opctlgen/`. Pass `-a name=value` for a required input it can't sample.
   On a mismatch, read the diff: it names the input set and the call. Fix the cause and rerun. Use `--force` only for a mismatch you
   understand and accept, and say why in the commit message.
4. **Resolve TODOs.** Each `TODO(opctl op convert)` marks a `RawRef` that reproduces op.yml exactly. Switch to the idiomatic form it names
   only when its condition holds for every value callers pass.
5. **Check callers.** Run `opctl op validate <op>`, then the op's callers and tests without editing them. Done when they pass and
   `git status` shows op.yml removed and op.go added, plus go.mod and .gitignore for a new module.

### Refusals

| Refused | Go rewrite, when there is one |
| --- | --- |
| `if` or loop reading a call result, `until`, `exists`/`notExists` | Split. |
| call `name` or `needs`; container `ports`; image `pullCreds` or `platform` | Split. |
| a string built from a call result or a secret in `cmd`, `workDir`, `name`, or the image | Bind the value with `EnvFrom` and build the string in the container's shell (`"$NAME"`). |
| a path or property of a call result, such as `$(out/sub)` or `$(meta.rev)` | Mount the whole dir and use the subpath inside the container, or declare the property as its own output. |
| a dir or file initializer | Commit the content as a file in the op directory and mount it with `opctl.OpFile` or `opctl.OpDir`. |
| `.go` files already in the op directory | Move them into a subdirectory. |
| Numeric loop ranges | Keep the loop in YAML. The current YAML interpreter rejects numeric ranges, and conversion preserves that behavior. |

A Go rewrite changes the graph, so convert first and rewrite in op.go afterwards; check the result with snapshots (step 4 below) and the
op's callers.

## Go to YAML

There's no converter for this direction: Go branches, helpers, and op stubs have no one-to-one op.yml form. Port by hand against
**snapshots**.

1. **Snapshot.** Run `opctl op plan <op> -a name=value ... > /tmp/<op>.<case>.yml` for the default inputs and for values that take each
   `if` branch and loop size in op.go. A snapshot is the graph for those inputs only. Done when every `if` and `for` in op.go lies between
   two snapshots that differ there.
2. **Write op.yml** from the default snapshot, then rewrite what plan baked in:
   - Interface: keep the snapshot's `name`, `description`, `inputs`, and `outputs`. Copy secret defaults from the op.go tags; snapshots omit them.
   - Plan-time reads: an `in.X.Value()` appears as the literal value. Replace it with `$(X)`.
   - Branches and loops: where snapshots differ, add an `if` with predicates over the input, or a `serialLoop` or `parallelLoop` over it.
   - Helpers: their calls appear inline. Keep them inline, or move a repeated block into its own op.
   - Generated names: rename each `opctl_<kind>_<n>` to a meaningful scope name. Names bound to outputs are already renamed.
   - Op stubs: they appear as plain `op:` calls by ref. Keep those calls.
3. **Remove the Go declaration.** Delete op.go, and also go.mod, go.sum, and .gitignore when only this op uses them. An op directory
   holding both op.yml and op.go is invalid.
4. **Compare.** For each saved case, run `opctl op plan` on the YAML op and diff it with the Go snapshot. Done when the only differences
   are scope names and every caller runs unchanged.
