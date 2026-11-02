---
name: converting-opctl-ops
description: Converts opctl ops between op.yml and Go. Use when porting a YAML op to a Go op, finishing a port that `opctl op convert --go` refused, or porting a Go op back to op.yml.
---

# Converting opctl ops

A conversion keeps the op's **interface** exactly: name, description, inputs, outputs, types, defaults, constraints, secrets, and
descriptions. Callers can't tell the op changed; only its declaration changes form.

## YAML to Go

1. **Preview.** From the directory holding `.opspec/`, run `opctl op convert --go <op> --dry-run`. Done when it prints either op.go with its module files
   (go.mod and .gitignore, plus go.sum with a published SDK, when the op gets its own module) or a list of **refusals**.
2. **Split refusals.** Skip when nothing was refused. Each refusal names an op.yml location and a construct the Go op SDK can't express.
   Move the smallest call containing it into a new YAML op beside this one, such as `.opspec/<op>-ports/op.yml`: give it an input for every
   `$(...)` the call reads and an output for every name later calls read, and replace the call in the original with an `op:` call to
   `../<op>-ports`. Check the **split**: `opctl op plan` of the original op before and after should differ only by the new `op:` call.
   Then return to step 1. Done when step 1 prints op.go. Some refusals have a rewrite instead; see [Refusals](#refusals).
3. **Convert.** Run `opctl op convert --go <op>`. It checks that the Go op plans the op.yml graph for the default inputs and generated
   samples, then writes the files, removes op.yml, and refreshes `opctlgen/`. Pass `-a name=value` for a required input it can't sample.
   On a mismatch, read the diff: it names the input set and the call. Fix the cause and rerun. Use `--force` only for a mismatch you
   understand and accept, and say why in the commit message.
4. **Resolve TODOs.** Each `TODO(opctl op convert)` marks a `RawRef` that reproduces op.yml exactly. Switch to the idiomatic form it names
   only when its condition holds for every value callers pass.
5. **Check callers.** Run `opctl op validate <op>`, then the op's callers and tests without editing them. Done when they pass and
   `git status` shows op.yml removed and op.go added, plus go.mod and .gitignore (and go.sum with a published SDK) when the op got its own module.

### Refusals

| Refused | Rewrite instead of a split |
| --- | --- |
| Go ops can't plan an empty op graph | Keep the op in YAML. |
| a root if/loop can plan an empty op graph (`a root if can plan an empty op graph`; `a root loop can plan an empty op graph`) | Split: move the `if` or loop into its own YAML op and call it unconditionally. |
| op refs can't contain whitespace or control characters | Rename the path and update its ref to a single token. Newlines in refs can inject generated Go source. |
| op ref must be a literal local or git ref; op ref is empty | Use a nonempty literal local path or git ref. Split if choosing the callee needs runtime data. |
| loop body reads and rebinds outer name `<name>` across iterations | Split: keep the loop and its carried value in a YAML op. |
| loop-body producer is referenced outside its loop; declared output is produced inside a loop | Split the loop and its consumers into a YAML op, with an explicit output contract. |
| if-branch producer is referenced outside its branch; declared output is produced inside an if branch and has no default | Keep the consumers in the same branch, or give a declared output a default for the skipped branch. |
| this position needs a plan-readable scalar; results, secrets and filesystem handles can't be read during plan | Split runtime predicates, loops, and image/name/workDir values. For cmd text, bind the whole value in `envVars` and let the container shell read `"$NAME"`. |
| until needs runtime data and can't be converted; exists and notExists predicates need runtime data | Split the loop or conditional call into a YAML op. |
| call name isn't supported by the Go op SDK; needs isn't supported by the Go op SDK | Split the named call and calls that depend on it into a YAML op. |
| container ports aren't supported by the Go op SDK; image pullCreds aren't supported by the Go op SDK; image platform isn't supported by the Go op SDK | Split the container call into a YAML op. |
| templates, paths, properties and items over produced values can't be converted; nested references over produced values can't be converted | Mount a whole produced dir and use its subpath inside the container. For a string or property, declare it as a separate output, or split. |
| dir initializers aren't supported by the Go op SDK; file initializers aren't supported by the Go op SDK | Commit the content as a file in the op directory and use a local path reference in the mount. |
| op directory already contains `<file>`; move Go sources into a subdirectory | Move the existing Go sources into a subdirectory before retrying. |
| op directory name isn't a valid Go module path element | Rename the directory to a valid Go module path element and update callers. |
| Numeric loop ranges (`numeric loop ranges aren't executable by the YAML interpreter`) | The op already fails at run time; YAML doesn't run numeric ranges. Use an array literal or array input as the range. |
| untyped loop item in a string position; add items: {type: string\|number\|boolean} to the array input | Add `items: {type: …}` to the array input. Choose `string`, `number`, or `boolean` to match its items. For a literal range, use one scalar type. |
| a loop needs an input or literal range; loop range must be a literal or a whole non-secret input; loop range can't read results or secrets; loop range must be an array or object | Use an array or object literal or a whole non-secret array or object input. Split when the range needs runtime data. |
| loop variable must bind a whole name | Use a whole name such as `$(item)` for each loop variable. |
| ordered predicates need number inputs or numeric literals; ordered predicates need numeric literals | Compare numbers. Use a number input or numeric literal; split when coercion needs runtime data. |
| constant predicate coercion error, such as `unable to coerce ... to number` | Correct the operands in op.yml. A constant predicate must be executable before it can be folded. |
| name `<name>` is produced by two parallel siblings (also `<path>`) | Give the outputs distinct names, or serialize the calls if one must replace the other. |
| cmd as one reference isn't supported; envVars as one reference isn't supported | Write `cmd` as an array and `envVars` as a mapping, or split when their shape needs runtime data. |
| mount needs a dir handle; mount needs a file handle; sockets must bind a whole socket input | Bind the matching whole handle. Use a dir or file path reference for an op-local mount. |
| typed container output has no output binding | Add the corresponding `outputs` binding for every declared `outputParams` entry. |
| callee has no input named `<name>`; callee has no output named `<name>` | Correct the call's bindings to match the callee's current interface. |
| declared output is unbound and has no default; output handle kind differs from its declaration | Bind the declared output to a handle of its declared kind, or supply a valid default when it can be unbound. |
| output binding must be null or a whole name | Use shorthand or a whole name such as `$(version)` for the output binding. |
| `parameter names must contain only letters, numbers, - or _, and can't start with opctl_; output binding starts with the reserved opctl_ prefix; references starting with opctl_ collide with generated names` | Rename the affected parameter or binding and update its references and callers. |
| null defaults and defaults of the wrong type can't be converted | Correct the default to match the parameter's declared kind. |
| unknown reference `<ref>` | Fix the name or declaration. Split when the reference depends on a runtime scope that Go can't express. |
| verification needs a value for input `<name>`; pass -a `<name>`=… | Pass a valid non-secret value with `-a name=value` and retry. |
| verification coverage incomplete: a true/false branch wasn't reached or loop length `<n>` wasn't reached | Supply inputs that reach the missing paths, or simplify unreachable control flow in op.yml. Use `--force` only after reviewing the reported gap. |
| verification loop at `<path>` exceeds 1000 iterations | Use a smaller verification range, or review the limit before accepting unverified output with `--force`. |
| op dir has a go.mod without the Go op SDK | Add the SDK requirement to the op's module, or remove that module file to let conversion create it. |

Each rewrite also works in op.yml. Make it there, check that callers still pass, then return to step 1, so verification compares against the op you ship.

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
