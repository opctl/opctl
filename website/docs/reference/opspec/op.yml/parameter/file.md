---
title: File parameter
---

An object defining a [file](../../types/file.md) parameter.

## Properties

### `default`

A literal file initialization used as the default value of the variable created by the parameter.

The default value can be a relative or absolute path. If it is an absolute path, the value is interpreted from the root of the op. If it is a relative path, the value is interpreted from the current working directory at the time the op is called. Relative path defaults are ignored when an op is called from another op as there is no current working directory.

### `isSecret`

A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's. 
