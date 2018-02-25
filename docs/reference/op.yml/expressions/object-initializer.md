Objects can be constructed from literal objects &/or expressions.

Expressions get replaced and [type coerced](type-coercion.md) (as required) on evaluation i.e. interpolation is supported.

ES2015 style shorthand property name syntax is supported.

```yaml
# literal object
myObject:
    prop1: value

# use property of scope object as object
$(obj.objProp)

# use json file at pkg root as object
$(/obj.json) 

# interpolate properties
myObject:
    prop1: string $(/someDir/file2.txt)
    $(prop2Name): $(someObj.someProp)
    prop3: [ sub, array, $(otherProp)]
    # Shorthand property name; equivalent to prop4: $(prop4)
    prop4:
```