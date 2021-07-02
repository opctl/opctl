---
title: Number Parameter [object]
---

An object defining a parameter which accepts a [number typed value](../../../types/number.md).

## Properties:
- may have:
  - [constraints](#constraints)
  - [default](#default)
  - [isSecret](#issecret)

### constraints
A [JSON Schema v4 [object]](https://tools.ietf.org/html/draft-wright-json-schema-00) defining constraints on the parameters value.

### default
A number to use as the value of the parameter when no argument is provided.

### isSecret
A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. 

## Example

This op echos a number parameter, defaulting to 9001, when run.

```yaml
name: example
description: an example op
inputs:
    example-input:
        number:
            default: 9001
run:
    container:
        image: { ref: 'alpine' }
        cmd: ['echo', $(example-input)]
```
