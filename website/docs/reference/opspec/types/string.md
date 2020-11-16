---
title: String
---

String typed values are a string i.e. an array of unicode characters

Strings...
- are immutable, i.e. assigning to an string results in a copy of the original string
- can be passed in/out of ops via [string parameters](../op-directory/op/parameter/string.md)
- can be initialized via [string initialization](#initialization)
- are coerced according to [string coercion](#coercion)

### Initialization
String typed values can be constructed from a literal or templated string.
 
A templated string is a string which includes one or more [variable-reference [string]](../op-directory/op/variable-reference.md).
At runtime, each reference gets evaluated and replaced with it's corresponding value.

#### Initialization Example (literal)

```yaml
i'm a string
```

#### Initialization Example (templated)
```yaml
# JSON representation of in scope object "someObject", replaces $(someObject) in the string
# contents of in scope file "someDir/file.txt" replaces $(someDir/file.txt) in the string
pre $(someObject) $(someDir/file.txt)
```

### Coercion
String typed values are coercible to:

- [boolean](boolean.md) (strings which are empty, all `"0"`, or (case insensitive) `"f"` or `"false"` coerce to false; all else coerce to true)
- [file](file.md)
- [number](number.md) (if value of string is numeric)
- [object](object.md) (if value of string is an object in JSON format)
