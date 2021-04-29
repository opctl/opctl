---
title: Object Parameter [object]
---

An object defining a parameter which accepts an [object typed value](../../../types/object.md).

## Properties:
- may have:
  - [constraints](#constraints)
  - [default](#default)
  - [isSecret](#issecret)

### constraints
A [JSON Schema v4 [object]](https://tools.ietf.org/html/draft-wright-json-schema-00) defining constraints on the parameters value.

### default
An object to use as the value of the parameter when no argument is provided.

### isSecret
A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. 

## Example
The op uses an object input, and echos the values of the object's keys.

```yaml
inputs:
    example-input:
        object:
            default:
                keyA: "valueA"
                keyB: ["subValueB1", "subValueB2", "subValueB3"]
                keyC:
                    subKeyC1: "subValueC1"
                    subKeyC2: "subKeyC2"
run:
    container:
        image: { ref: 'apteno/alpine-jq' }
        cmd:
            - sh
            - -ce
            - echo $(example-input.keyA) && echo $(example-input.keyB[0]) && echo $(example-input.keyC.subKeyC1)
```

The expected output is:

```sh
valueA
subValueB1
subValueC1
```
