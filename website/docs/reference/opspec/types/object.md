---
title: Object
---

Object values are a container for string indexed values (referred to as properties). They behave as key-value maps, where the keys are strings. Objects are immutable: assigning to an object produces a copy of the original object.

## Initialization

Objects are initialized with standard yaml object syntax. [Variable references](../op.yml/variable-reference.md) within the object's keys will be evaluated as strings and replaced when initialized. References within values will be evaluated at initialization and will retain their types. Keys with no (`null`) values are treated as shorthand: the value will be treated as a variable reference to the key name.

```yaml
myObject:
  prop1: string $(./someDir/file2.txt) # will be interpolated as a string, replacing the reference with the file content
  $(prop2Name): $(someObject.someProperty) # the key will take the value of prop2Name as a string
  prop3: [sub, array, 2] # { "prop3": ["sub", "array", 2] }
  prop4: # This shorthand is equivalent to `prop4: $(prop4)`
```

## Coercion

Object values are coercible to:

- [boolean](boolean.md): `null` or empty objects coerce to `false`, all else coerces to `true`
- [file](file.md): The object will be serialized to JSON
- [dir](dir.md): The object will be treated as a [`dir` initializer](dir#literal-initialization)
- [string](string.md): The object will be serialized to JSON

```yaml
name: objAsString
inputs:
  obj:
    object:
      default:
        prop1: prop1Value
        prop2: [item1]
run:
  container:
    image: { ref: alpine }
    cmd:
      - echo
      - $(obj)
```
