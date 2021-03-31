---
title: String
---

String values are a string (an array of unicode characters)

## Initialization

Strings can be created with a yaml string literal. They can contain [variable references](../op.yml/variable-reference.md) that will be evaluated, coerced to a string, and replaced during initialization.

```yaml
inputs:
  myInput:
    number:
      default: 3 # 3!
run:
  op:
    ref: ../op
    inputs:
      input1: A simple string
      input2: "Another simple string"
      input3: |
        A more complex string with multiple lines
        and a variable reference: $(myInput)
```

## Coercion

String values are coercible to:

- [boolean](boolean.md): Empty, all `"0"`, or (case insensitive) `"f"` or `"false"` strings coerce to `false`, all else coerce to `true`
- [file](file.md)
- [number](number.md): if the value is numeric
- [object](object.md): if the value is a valid JSON object, will be parsed as JSON
