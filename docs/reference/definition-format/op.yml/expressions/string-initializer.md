Strings can be constructed from literal text &/or expressions.

Expressions get replaced and [type coerced](../type-coercion.md) (as required) on evaluation i.e. interpolation is supported. 

## Examples

### Literal

```yaml
i'm a string
```

### Interpolated
given:
- someObject
  - is in scope
  - is object

```yaml
# $(someObject) replaced w/ JSON representation of someObject
# $(dir/file.txt) replaced w/ contents of file.txt
pre $(someObject) $(dir/file.txt)
```