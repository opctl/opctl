Interpolation is supported where noted.

# String interpolation

Used for string templating.

Replaces value expressions contained in strings w/ their value.

> values will be [type coerced](type-coercion.md) to string if
> applicable.

# Number interpolation

Used for number templating.

Replaces value expressions contained in numbers w/ their values.

> values will be [type coerced](type-coercion.md) to number if
> applicable.

# Examples

#### String w/ interpolation

```yaml
name: stringInterpolation
inputs:
  greeting:
    string:
      constraints: { enum: [hello, hi, welcome]}
      default: hello
  aNumber:
    number:
      constraints: { format: integer }
      default: 2
  anObject:
    object:
      default:
        prop1: prop1Value
run:
  container:
    image: { ref: alpine }
    cmd:
    - echo
    - $(greeting) $(aNumber), look at this JSON $(anObject)
```

