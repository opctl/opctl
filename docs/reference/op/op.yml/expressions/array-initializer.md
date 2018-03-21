Arrays can be constructed from literal arrays &/or expressions.

Expressions get replaced and [type coerced](../type-coercion.md) (as required) on evaluation i.e. interpolation is supported. 

```yaml
# literal array
- item1
- item2

# use property of scope object as array
$(obj.arrayProp)

# use json file embedded in op as array
$(/array.json) 

# interpolate items
- string $(/someDir/file2.txt)
- $(someObj.someProp)
- [ sub, array, $(otherProp)]
```