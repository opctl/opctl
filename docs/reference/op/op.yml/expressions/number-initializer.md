Numbers can be constructed from literal numbers &/or expressions.

Expressions get replaced and [type coerced](../type-coercion.md) (as required) on evaluation i.e. interpolation is supported. 

```yaml
# interpolate w/ number property of scope object
222$(obj.propCoercibleToNum)3

# interpolate w/ number file embedded in op
44e$(/fileCoercibleToNum)
```