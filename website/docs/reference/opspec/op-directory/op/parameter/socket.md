---
title: Socket Parameter [object]
---

An object defining a parameter which accepts a [socket typed value](../../../../types/socket.md).

## Socket Properties:
- must have:
  - [description](#description)
- may have:
  - [isSecret](#issecret)

### description
A [markdown [string]](../markdown.md) defining a human friendly description of the parameter.

### isSecret
A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. 
