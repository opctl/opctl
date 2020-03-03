---
title: Boolean
---

Boolean typed values are a boolean i.e. `true` or `false`.

Booleans...
- are immutable, i.e. assigning to a boolean results in a copy of the original boolean
- can be passed in/out of ops via [boolean parameters](../structure/op-directory/op/parameter/boolean.md)
- can be initialized via [boolean initialization](#initialization)
- are coerced according to [boolean coercion](#coercion)

### Initialization
Boolean typed values can be constructed from a literal boolean.

### Coercion
Boolean typed values are coercible to:

- [file](file.md) (will be serialized to JSON)
- [string](string.md) (will be serialized to JSON)
