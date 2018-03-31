Object properties can be referenced via `.propertyName` syntax.

## Examples

### From scope
given:
- `someObject`
  - is in scope
  - is type coercible to object
  - contains property `someProperty`

```yaml
$(someObject.someProperty)
```