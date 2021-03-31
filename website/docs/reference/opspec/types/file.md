---
title: File
---

File values are a filesystem file entry. Files are passed by reference, not copied, which means changing the contents or value of a file will cause it to change in each variable that holds it.

## Initialization

File values can be constructed literally or through a [variable reference](../op.yml/variable-reference.md).

Literal values will be interpreted as [strings](string).

```yaml
myLiteralFile: |
  This will be the contents of the file
myTemplatedLiteralFile: |
  This will be the contents of the file, which supports templating: $(myVar) 
myReferencedFile: $(myDir/file)
```

## Coercion

File values are coercible to:

- [boolean](boolean.md): files which are empty, all `"0"`, or (case insensitive) `"f"` or `"false"` coerce to `false`, all else coerce to `true`
- [array](array.md): if the content of the file is a valid JSON array, will be parsed as JSON
- [number](number.md): if the content of the file is numeric
- [object](object.md): if the content of the file is a valid JSON object, will be parsed as JSON
- [string](string.md): content of the file
