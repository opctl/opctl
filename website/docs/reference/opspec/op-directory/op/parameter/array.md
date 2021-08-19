---
title: Array Parameter [object]
---

An object defining a parameter which accepts an [array typed value](../../../types/array.md).

## Properties:
- may have:
  - [constraints](#constraints)
  - [default](#default)
  - [isSecret](#issecret)

### constraints
A [JSON Schema v4 [object]](https://tools.ietf.org/html/draft-wright-json-schema-00) defining constraints on the parameters value.

### default
An [array initializer](../../../types/array.md#initialization) to use as the value of the parameter when no argument is provided.

### isSecret
A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. 

## Example

This is an example op that takes in an array parameter and prints the values in the array

```yaml
name: example
description: an example op
inputs:
  example-input:
    array:
      default:
        - "hello"
        - "world"
        - "!"
run:
  container:
    image: { ref: 'alpine' }
    cmd: ["echo", $(example-input)]
```
