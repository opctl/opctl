---
title: File Parameter [object]
---

An object defining a parameter which accepts a [file typed value](../../../types/file.md).

## Properties:
- may have:
  - [default](#default)
  - [isSecret](#issecret)

### default
A relative or absolute path string to use as the default value of the parameter when no argument is provided.

If the value is...
- an absolute path, the value is interpreted from the root of the op.
- a relative path, the value is interpreted from the current working directory at the time the op is called.
  > relative path defaults are ignored when an op is called from an op as there is no current working directory.

### isSecret
A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. 

## Example
This is an example op that echos the contents of a file input

```yaml
name: example
description: an example op
inputs:
    example-input:
        file:
            default: ./optest
run:
    container:
        image: { ref: 'alpine' }
        cmd: ['echo', $(example-input)]
```

To use this example, create a file called `optest` in the same directory the op is called from.  Whatever content is in the `optest` file will be echoed.

For example, the following commands will result in "hello world" being echoed by the example op.
```sh
echo hello world > optest
opctl run example
```
