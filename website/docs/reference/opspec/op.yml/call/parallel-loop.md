---
title: Parallel loop call
---

A parallel loop call is an object that runs a call multiple times in parallel (concurrently with no defined order). Examples of use cases include running multiple identical calls (for load testing, for example) or running multiple parameterized calls (such as an identical build for different target architectures).

## Properties

### `range`

_required_

A [rangeable value](../rangeable-value.md) to loop over.

### `run`

_required_

A [call](index.md) that will be run each iteration of the loop. [Loop variable](#vars) will be in scope in the definition of the call.

### `vars`

A [loop variables declaration](../loop-vars.md) that exposes values from the `range` to the call.

## Example

The following op will emit:

```
firstKey firstValue 0
secondKey secondValue 1
```

```yaml
description: An example of a parallel loop
run:
  parallelLoop:
    range:
      firstKey: firstValue
      secondKey: secondValue
    vars:
      index: $(index)
      key: $(key)
      value: $(val)
    run:
      container:
        image: { ref: alpine }
        cmd: ["echo", $(key), $(val), $(index)]
```
