Strings can be constructed from literal text &/or expressions.

Expressions get replaced and [type coerced](type-coercion.md) (as required) on evaluation i.e. interpolation is supported. 

```yaml
"here's json: $(obj)"
"here's contents of a file: (dir/file.txt)"
```