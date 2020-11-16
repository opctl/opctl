---
title: Variable Reference [string]
---
A string referencing the location of/for a value in the form of `$(REFERENCE)` where `REFERENCE` MUST start with an [identifier [string]](identifier.md) and MAY end with one or more:
- [array item references](../../types/array.md#item-referencing)
- [object property references](../../types/object.md#property-referencing)
- [dir entry references](../../types/dir.md#entry-referencing)

References can be used to either define or access values in the current scope. 

When an op starts, it's initial scope includes:

- `./` equal to the current op directory i.e. the current `op.yml` can be accessed via `$(./op.yml)`.
- `../` equal to the parent of the current op directory i.e. the current `op.yml` can be accessed via `$(../op.yml)`.
- any defined inputs


> note: variable references can be escaped by prefixing the [would be] variable reference with `\` i.e. `\\$(wouldBeVariableReference)` would not be treated as a variable reference. 