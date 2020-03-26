---
title: File Parameter [object]
---

An object defining a parameter which accepts a [file typed value](../../../../types/file.md).

## Properties:
- must have:
  - [description](#description)
- may have:
  - [default](#default)
  - [isSecret](#issecret)

### default
A relative or absolute path string to use as the default value of the parameter when no argument is provided.

If the value is...
- an absolute path, the value is interpreted from the root of the op.
- a relative path, the value is interpreted from the current working directory at the time the op is called.
  > relative path defaults are ignored when an op is called from an op as there is no current working directory.

### description
A [markdown [string]](../markdown.md) defining a human friendly description of the parameter.

### isSecret
A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. 
