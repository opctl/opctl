---
title: File Parameter [object]
---

An object defining a parameter which accepts a [file typed value](../../../types/file.md).

## Properties:
- may have:
  - [default](#default)
  - [isSecret](#issecret)

### default
A [file initializer](../../../types/file.md#initialization) to use as the value of the parameter when no argument is provided.

If the value is a relative path it will be resolved from the current working directory of the caller. If no current working directory exists, such as when the caller is an op or web UI, the default will be ignored.

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
      default: $(./optest)
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
