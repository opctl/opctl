---
title: Boolean
---

Boolean values are `true` or `false`. Booleans are immutable: assigning a boolean variable to another variable copies the value.

## Initialization

Booleans are created with [yaml boolean literal syntax](https://yaml.org/type/bool.html).

```yaml
run:
  op:
    ref: ../op
    inputs:
      input1: true
```

## Coercion

Boolean values are coercible to:

- [file](file.md): The value will be serialized to JSON
- [string](string.md) The value will be serialized to JSON
