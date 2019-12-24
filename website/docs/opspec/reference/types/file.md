---
title: File
---

File typed values are a filesystem file entry.

Files...
- are mutable, i.e. making changes to a file results in the file being changed everywhere it's referenced
- can be passed in/out of ops via [file parameters](../structure/op-directory/op/parameter/file)
- can be initialized via [file initialization](#initialization)
- are coerced according to [file coercion](#coercion)

### Initialization
File typed values can be constructed from a literal string or templated string; see [string initialization](string#initialization).

### Coercion
File typed values are coercible to:

- [boolean](#boolean) (files which are empty, all `"0"`, or (case insensitive) `"f"` or `"false"` coerce to false; all else coerce to true)
- [array](array) (if value of file is an array in JSON format)
- [number](number) (if value of file is numeric)
- [object](object) (if value of file is an object in JSON format)
- [string](string)