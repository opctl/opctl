---
sidebar_label: Index
title: Parameter [object]
---

An object defining a parameter of an operation; i.e. a value that is passed into or out of it's scope.

## Properties
- must have exactly one of
  - [array](array.md)
  - [boolean](boolean.md)
  - [dir](dir.md)
  - [file](file.md)
  - [number](number.md)
  - [object](object.md)
  - [socket](socket.md)
  - [string](string.md)

## Example
This is an example op with one of each type of input parameter.

```yaml
name: example
description: an example op
inputs:
  example-input-array:
    array:
      description: "An example array input parameter"
      default:
        - "a"
        - "b"
        - "c"
  example-input-boolean:
    boolean:
      description: "An example boolean input parameter"
      default: false
  example-input-dir:
    dir:
      description: "An example input dir parameter"
      default: .
  example-input-file:
    file:
      description: "An example input file parameter"
      default: .opspec/example/op.yml
  example-input-number:
    number:
      description: "An example input number"
      default: 1
  example-input-object:
    object:
      description: "An example input object"
      default:
        example-object-field-1: "field 1"
        example-object-field-2: "field 2"
  example-input-socket:
    socket:
      description: "An example input socket"
  example-input-string:
    string:
      description: "An example input string"
      default: "a default value"
run:
  container:
    image: { ref: 'alpine' }
    cmd: ['echo', $(example-input-array)]
```
