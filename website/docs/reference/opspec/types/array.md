---
title: Array
---

Array typed values are a container for numerically indexed values (referred to as items).

Arrays...
- are immutable, i.e. assigning to an array results in a copy of the original array
- can be passed in/out of ops via [array parameters](../op-directory/op/parameter/array.md)
- can be initialized via [array initialization](#initialization)
- items can be referenced via [array item referencing](#item-referencing)
- are coerced according to [array coercion](#coercion)

### Initialization
Array typed values can be constructed from a literal or templated array.
 
A templated array is an array which includes one or more [variable-reference [string]](../op-directory/op/variable-reference.md).
At runtime, each reference gets evaluated and replaced with it's corresponding value. 

#### Initialization Example (literal)
```yaml
- item1
- item2
```

#### Initialization Example (templated)
given:
- `someDir/file2.txt` is embedded in op
- `someObject` 
  - is in scope
  - is type coercible to object
  - has property `someProperty`

```yaml
- string $(./someDir/file2.txt)
- $(someObject.someProperty)
- [ sub, array, 2]
```

### Item Referencing
Array items can be referenced via `$(ARRAY[index])` syntax, where `index` is the zero based index of the item. 
If `index` is negative, indexing will take place from the end of the array.

#### Item Referencing Example (first item)
given:
- someArray
  - is in scope
  - is type coercible to array
  - has at least one item

```yaml
$(someArray[0])
```

#### Item Referencing Example (last item)
given:
- someArray
  - is in scope
  - is type coercible to array
  - has at least one item

```yaml
$(someArray[-1])
```

### Coercion
Array typed values are coercible to:

- [boolean](boolean.md) (arrays which are null or empty coerce to false; all else coerce to true)
- [file](file.md) (will be serialized to JSON)
- [string](string.md) (will be serialized to JSON)