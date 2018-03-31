Objects can be constructed from literal objects &/or expressions.

Expressions get replaced and [type coerced](../type-coercion.md) (as required) on evaluation i.e. interpolation is supported.

ES2015 style shorthand property name syntax is supported.

## Examples

### Literal

```yaml
myObject:
    prop1: value
```

### Object property ref
given:
- `someObject`
  - is in scope
  - is type coercible to object
- `objectProperty`
  - is a property of `someObject`
  - is type coercible to array

```yaml
$(someObject.objectProperty)
```

### Embedded dir entry ref
given:
- `/object.json`
  - is embedded in op
  - is type coercible to object

```yaml
$(/object.json) 
```

### Interpolated
given:
- `/someDir/file2.txt` is embedded in op
- `prop2Name` is in scope
- `someObject`
  - is in scope
  - is type coercible to object
  - has property `someProperty`
- `prop4` is in scope

```yaml
# interpolate properties
myObject:
    prop1: string $(/someDir/file2.txt)
    $(prop2Name): $(someObject.someProperty)
    prop3: [ sub, array, 2]
    # Shorthand property name; equivalent to prop4: $(prop4)
    prop4:
```