---
title: Boolean Parameter [object]
---

An object defining a parameter which accepts a [boolean typed value](../../../types/boolean.md).

## Boolean Properties:
- may have:
  - [default](#default)
  - [isSecret](#issecret)

### default
A boolean to use as the value of the parameter when no argument is provided.

## Example
This is an example op that echos a boolean input.

```yaml
name: example
description: an example op
inputs:
    example-input:
        boolean:
            default: false
run:
    container:
        image: { ref: 'alpine' }
        cmd: ["echo", $(example-input)]
```
