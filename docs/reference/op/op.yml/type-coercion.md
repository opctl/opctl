Type coercion takes place automatically when necessary/possible.

# Array coercion

Array typed values are coercible to:

- file (will be serialized to JSON)
- string (will be serialized to JSON)

# File coercion

File typed values are coercible to:

- array (if value of file is an array in JSON format)
- number (if value of file is numeric)
- object (if value of file is an object in JSON format)
- string

# Number coercion

Number typed values are coercible to:

- file
- string

# Object coercion

Object typed values are coercible to:

- file (will be serialized to JSON)
- string (will be serialized to JSON)

# String coercion

String typed values are coercible to:

- file
- number (if value of string is numeric)
- object (if value of string is an object in JSON format)

# Examples

## Object to string

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

## Number to file

```yaml
name: numAsFile
run:
  container:
    image: { ref: alpine }
    cmd:
    - sh
    - -ce
    - cat /numCoercedToFile
    files:
      /numCoercedToFile: 2.2
```

