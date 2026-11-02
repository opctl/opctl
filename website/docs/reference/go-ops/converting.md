---
sidebar_position: 6
---

# Converting ops

`opctl op convert --go` changes a local YAML op into a Go op while preserving its interface and op graph. Existing callers keep the same ref, inputs, and outputs.

```sh
opctl op convert --go dev --dry-run
opctl op convert --go dev
```

Run these commands from the directory holding `.opspec/`, or pass a local path. Git refs aren't accepted by `convert`.

## Preview and verify

A dry run analyzes the op file, renders Go, and verifies it through the node. It prints each proposed file under a `// ==> <path> <==` line and leaves the op unchanged.
The report goes to stderr.

Conversion compares the Go op's interface and plans against the op file for the default inputs and up to 16 generated input sets. It uses samples for required
public value inputs, branch alternatives, and loop sizes. Supply a value when the generated samples can't satisfy an input's constraints:

```sh
opctl op convert --go deploy -a region=us-east-1 --arg-file inputs.yml
```

Only non-secret string, number, boolean, array, and object inputs participate in verification. Arguments for secret, dir, file, or socket inputs are ignored,
with a note naming the input. Verification doesn't execute the op's calls. It needs a running node and the Go bootstrap toolchain.

A mismatch reports the input set and the first differing graph or interface field. Nothing is written. Inspect the generated files with `--dry-run` and fix the cause.
`--force` writes an unverified Go op with a warning when verification fails or can't run. It never overrides an analysis refusal. A dry run still fails verification
when `--force` is present.

Verification checks generated samples, so run the op's existing tests and callers after conversion.

## Files and modules

On a successful conversion, opctl writes `op.go`, removes `op.yml`, and refreshes the module's git-ignored `opctlgen/` editor stubs. A stub refresh failure is a warning.
Bootstrap still generates the stubs it needs when the op runs.

If the enclosing module already requires `github.com/opctl/opctl/sdks/go/opspec/sdk`, conversion reuses it. Otherwise the op receives its own `go.mod` and `.gitignore`,
plus `go.sum` when this opctl build supplies a published SDK. An existing `.gitignore` gets `/opctlgen/` appended when needed. Conversion doesn't edit an enclosing
application's `go.mod`.

An op with its own `go.mod` that doesn't require the SDK must add that requirement or remove the module file before converting. Development builds without a published
SDK need a local `replace` directive in a prepared module before verification can run.

Module files are written first, then `op.go` is installed before `op.yml` is removed. An interruption can leave both declarations. Keep the intended declaration and
remove the other before running the op.

## Generated Go

Inputs and outputs keep their names, types, defaults, constraints, descriptions, and secret flags. The op's name and description remain the same. Deprecated nested
parameter descriptions move to the parameter. The report identifies dropped op versions and call descriptions retained as comments.

Input-only branches and loops become Go `if` and `for` statements. Arrays and objects retain their schemas in tags. A scalar item schema can give an array loop input
a concrete Go element type.

Op calls use [typed op stubs](op-stubs.md) when every bound value and consumed output has the exact required Go type. Other calls use `CallOp`, with a comment
explaining the type mismatch.

Some input references need `RawRef` to preserve YAML behavior. Each has a `TODO(opctl op convert)` explaining the idiomatic form and the condition under which it is
equivalent. Resolve it only when that condition holds for every value callers pass.

## Refusals

Analysis lists every refusal with its op-file line, column, and field. It writes nothing, including with `--force`.

Common refusals include:

- Runtime-dependent branches or loops, `until`, and `exists` or `notExists` predicates.
- Call `name` or `needs`, container `ports`, and image `pullCreds` or `platform`.
- Templates, paths, properties, or items over call results, and secrets in positions that need literal strings during plan.
- Dir or file initializers, incompatible defaults, unbound outputs, and existing Go sources in the op directory.
- Numeric loop ranges. The YAML interpreter rejects them; conversion preserves that behavior.

Move the smallest unsupported call into a YAML op with explicit inputs and outputs, replace it with an op call, and retry conversion. The
[converting-opctl-ops skill](https://github.com/opctl/opctl/tree/main/.agents/skills/converting-opctl-ops) explains the split and possible Go rewrites.

## Plan snapshots

`opctl op plan` prints an op-file-shaped snapshot of the graph a Go or YAML op plans for one input set:

```sh
opctl op plan deploy -a region=us-east-1 > deploy.east.yml
opctl op plan deploy --arg-file inputs.yml > deploy.default.yml
```

It accepts the same op refs as `opctl run`, including git refs. It never executes the op's calls. Go ops still run their describe and plan bootstrap phases.
YAML ops expand input-only branches and array or object loops. Runtime-dependent constructs and loops above 1,000 iterations remain in the snapshot.

The header records the public value inputs used for the snapshot. Reference-containing defaults remain in the interface rather than appearing as resolved header
values. Secret, dir, file, and socket inputs stay references. Secret defaults are replaced with a comment, and explicitly supplied secret values aren't sent to plan.

A snapshot is for inspection and debugging. Different inputs can produce a different graph, and Go branches, helpers, and op stubs don't appear as Go constructs.
There is no automatic Go-to-YAML conversion. Preserve the interface and use snapshots for every branch and loop size while porting by hand.

Agents can install the conversion workflow with:

```sh
npx skills add opctl/opctl
```
