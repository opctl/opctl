Numbers can be constructed from literal numbers &/or expressions.

Expressions get replaced and [type coerced](../type-coercion.md) (as required) on evaluation i.e. interpolation is supported. 

## Examples

### Literal
```yaml
2
```

### Interpolated
given:
- `someNumber`
  - is in scope
   - is type coercible to number

```yaml
# $(someNumber) replaced w/ someNumber
222$(someNumber)3e10
```