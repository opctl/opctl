A single item of an array can be referenced via `[index]` syntax, where `index` is the zero based index of the item. 
If `index` is negative, indexing will take place from the end of the array. 

## Examples

### First item
given:
- someArray
  - is in scope
  - is type coercible to array
  - has at least one item

```yaml
$(someArray[0])
```

### Last item
given:
- someArray
  - is in scope
  - is type coercible to array
  - has at least one item

```yaml
$(someArray[-1])
```