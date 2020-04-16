---
title: Object
---

Object typed values are a container for string indexed values (referred to as properties).

Objects...
- are immutable, i.e. assigning to an object results in a copy of the original object
- can be passed in/out of ops via [object parameters](../op-directory/op/parameter/object.md)
- can be initialized via [object initialization](#initialization)
- properties can be referenced via [object property referencing](#property-referencing)
- are coerced according to [object coercion](#coercion)

### Initialization
Object typed values can be constructed from a literal or templated object.
 
A templated object is an object which includes one or more value reference.
At runtime, each reference gets evaluated and replaced with it's corresponding value.

#### Initialization Example (literal)

```yaml
myObject:
    prop1: value
```

#### Initialization Example (templated)
given:
- `/someDir/file2.txt` is embedded in op
- `prop2Name` is in scope
- `someObject`
  - is in scope
  - is type coercible to object
  - has property `someProperty`
- `prop4` is in scope

```yaml
# interpolate properties
myObject:
    prop1: string $(/someDir/file2.txt)
    $(prop2Name): $(someObject.someProperty)
    prop3: [ sub, array, 2]
    # Shorthand property name; equivalent to prop4: $(prop4)
    prop4:
```

### Property Referencing
Object properties can be referenced via `$(OBJECT.PROPERTY)` or `$(OBJECT[PROPERTY])` syntax.

#### Property Referencing Example (from scope)
given:
- `someObject`
  - is in scope
  - is type coercible to object
  - contains property `someProperty`

```yaml
$(someObject.someProperty)
```

### Coercion
Object typed values are coercible to:

- [boolean](boolean.md) (objects which are null or empty coerce to false; all else coerce to true)
- [file](file.md) (will be serialized to JSON)
- [string](string.md) (will be serialized to JSON)

#### Coercion Example (object to string)
```yaml
name: objAsString
inputs:
  obj:
    object:
      default:
        prop1: prop1Value
        prop2: [ item1 ]
run:
  container:
    image: { ref: alpine }
    cmd:
    - echo
    - $(obj)
```