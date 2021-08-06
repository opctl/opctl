---
title: Dir Parameter [object]
---

An object defining a parameter which accepts a [dir typed value](../../../types/dir.md).

## Properties:
- may have:
  - [default](#default)
  - [isSecret](#issecret)

### default
A [dir initializer](../../../types/dir.md#initialization) to use as the value of the parameter when no argument is provided.

If the value is a relative path it will be resolved from the current working directory of the caller. If no current working directory exists, such as when the caller is an op or web UI, the default will be ignored.

### isSecret
A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. 
