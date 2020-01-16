---
title: Reference [string]
---
A string used to reference the source of or destination for a value.

The format is `$(LOCATION)` where `LOCATION` MUST be one of:

- an absolute path to a file or directory included (a.k.a embedded) in the current operation e.g. the current op.yml can be referenced via `$(/op.yml)`.
- a [variable name](variable-name.md) representing a scoped variable.