Arrays can be constructed from literal arrays &/or expressions.

Expressions get replaced and [type coerced](../type-coercion.md) (as required) on evaluation i.e. interpolation is supported. 

## Examples

### Literal
```yaml
- item1
- item2
```

### Object property ref
given:
- `someObject`
  - is in scope
  - is type coercible to object
- `arrayProperty`
  - is a property of `someObject`
  - is type coercible to array

```yaml
$(someObject.arrayProperty)
```

### Embedded dir entry ref
given:
- `/array.json`
  - is embedded in op
  - is type coercible to array

```yaml
$(/array.json) 
```

### Interpolated
given:
- `/someDir/file2.txt` is embedded in op
- `someObject` 
  - is in scope
  - is type coercible to object
  - has property `someProperty`

```yaml
- string $(/someDir/file2.txt)
- $(someObject.someProperty)
- [ sub, array, 2]
```