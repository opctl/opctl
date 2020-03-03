---
title: Number Parameter [object]
---

An object defining a parameter which accepts a [number typed value](../../../../types/number.md).

## Properties:
- must have:
  - [description](#description)
- may have:
  - [constraints](#constraints)
  - [default](#default)
  - [isSecret](#issecret)

### constraints
A [JSON Schema v4](https://tools.ietf.org/html/draft-wright-json-schema-00) object defining constraints on the parameters value.

### default
A number to use as the value of the parameter when no argument is provided.

### description
A [markdown](../markdown.md) string defining a human friendly description of the parameter.

### isSecret
A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. 
