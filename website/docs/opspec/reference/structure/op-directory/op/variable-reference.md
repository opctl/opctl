---
title: Variable Reference [string]
---
A string used to reference a variable in the form of `$(VARIABLE)` where `VARIABLE` MUST be a valid [variable name [string]](variable-name.md).

Variable references can be used to either define or access variables in the current scope. 

When an op starts, it's initial scope includes:

- `/` with a value of the current op directory i.e. the current op's `op.yml` can be accessed via `$(/op.yml)`.
- any defined inputs


> note: variable references can be escaped by prefixing the [would be] variable reference with `\` i.e. `\\$(wouldBeVariableReference)` would not be treated as a variable reference. 