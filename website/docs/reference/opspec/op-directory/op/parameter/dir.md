---
title: Dir Parameter [object]
---

An object defining a parameter which accepts a [dir typed value](../../../types/dir.md).

## Properties:
- may have:
  - [default](#default)
  - [isSecret](#issecret)

### default
A relative or absolute path to use as the default value of the parameter when no argument is provided.

If the value is...
- an absolute path, the value is interpreted from the root of the op.
- a relative path, the value is interpreted from the current working directory at the time the op is called.
  > relative path defaults are ignored when an op is called from an op as there is no current working directory.

### isSecret
A boolean indicating if the value of the parameter is secret. This will cause it to be hidden in UI's for example. 

## Example
This is an example op that lists the contents of a dir input.  Note that the dir input is mounted in the container as `/mounted-example-input`

```yaml
name: example
description: an example op
inputs:
    example-input:
        dir:
            default: .
run:
    container:
        image: { ref: 'alpine' }
        cmd: ["ls", "-la", /mounted-example-input]
        dirs:
            /mounted-example-input: $(example-input)
```

The expected output is the same as running `ls -la` from the directory the op is called from
