---
title: Number
---

Numbers are numerical values. Opspecs do not define floating or integer semantics, but stores numbers internally as `float64` values. Numbers are immutable: assigning a number variable to another variable copies the value.

## Initialization

Numbers are initalized with normal yaml number syntax.

Because strings can be coerced to numbers, strings that include [variable references](../op.yml/variable-reference.md) will be interpreted and coerced to a number.

```yaml
inputs:
  myInput:
    number:
      default: 3 # 3!
run:
  op:
    ref: ../op
    inputs:
      input1: 12$(myInput)4e10 # will be interpreted as 12340000000000
```

## Coercion

Number values are coercible to:

- [boolean](boolean.md): `0` is `false`, all else is `true`
- [file](file.md)
- [string](string.md)

<!-- TODO: This example doesn't work -->

```yaml
name: numAsFile
run:
  container:
    image: { ref: alpine }
    cmd:
      - cat
      - /numCoercedToFile
    files:
      /numCoercedToFile: 2.2
```
