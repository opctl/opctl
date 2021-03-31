---
title: Array
---

Array values are contain numerically indexed values ("items"). Arrays are immutable: assigning to an array produces a copy of the original array.

## Initialization

Arrays are initialized with standard yaml sequence syntax. [Variable references](../op.yml/variable-reference.md) within the array will be evaluated and replaced with the value when initialized.

```yaml
inputs:
  myString:
    string:
      default: "hello world"
run:
  op:
    ref: ../op
    inputs:
      input1:
        - item1 # string literal, "item1"
        - $(myString) # will be unwrapped to "hello world"
        - [1, 2, 3] # a sub-array
```

## Coercion

Array values are coercible to:

- [boolean](boolean.md): `null` or empty arrays coerce to `false`, all else coerces to `true`
- [file](file.md): The array will be serialized to JSON
- [string](string.md): The array will be serialized to JSON
