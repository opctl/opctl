---
title: String Parameter [object]
---

An object defining a parameter which accepts a [string typed value](../../../types/string.md).

## Properties:
- must have:
  - [description](#description)
- may have:
  - [constraints](#constraints)
  - [default](#default)
  - [isSecret](#issecret)

### constraints
A [JSON Schema v4 [object]](https://tools.ietf.org/html/draft-wright-json-schema-00) defining constraints on the parameters value.

#### default
A string to use as the value of the parameter when no argument is provided.

#### description
A [markdown [string]](../markdown.md) defining a human friendly description of the parameter.

#### isSecret
A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example.

## Example

This is an example op that uses a string input, with a default value

```yaml
name: example
description: an example op
inputs:
  example-input:
    string:
      default: "a default value"
run:
  container:
    image: { ref: 'alpine' }
    cmd: ['echo', $(example-input)]
```

Using the example op:
```shell-script
opctl run example
```
The expected output is the op running and echoing "a default value".

Using the example op while overriding the default value:
```shell-script
opctl run -a example-input="hello world" example
```
The expected output is the op running and echoing "hello world".
