---
title: String
---

String typed values are a string i.e. an array of unicode characters

Strings...
- are immutable, i.e. assigning to an string results in a copy of the original string
- can be passed in/out of ops via [string parameters](../structure/op-directory/op/parameter/string)
- can be initialized via [string initialization](#initialization)
- are coerced according to [string coercion](#coercion)

### Initialization
String typed values can be constructed from a literal or templated object.
 
A templated string is a string which includes one or more value reference.
At runtime, each reference gets evaluated and replaced with it's corresponding value.

#### Initialization Example (literal)

```yaml
i'm a string
```

#### Initialization Example (templated)
given:
- someObject
  - is in scope
  - is object

```yaml
# $(someObject) replaced w/ JSON representation of someObject
# $(dir/file.txt) replaced w/ contents of file.txt
pre $(someObject) $(dir/file.txt)
```

### Coercion
String typed values are coercible to:

- [boolean](#boolean) (strings which are empty, all `"0"`, or (case insensitive) `"f"` or `"false"` coerce to false; all else coerce to true)
- [file](#file)
- [number](#number) (if value of string is numeric)
- [object](#object) (if value of string is an object in JSON format)
