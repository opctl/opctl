---
title: Serial loop call
---

A serial loop call is an object that runs a call multiple times one after another. Examples of use cases include wait for a process to start or restarting a process if it crashes.

## Properties

Serial loops must define a lifetime using one of `range` or `until`.

### `range`

_required if no `until` specified_

A [rangeable value](../rangeable-value.md) to loop over.

### `run`

_required_

A [call](../call/index.md) that will be run each iteration of the loop. [Loop variable](#vars) will be in scope in the definition of the call.

### `until`

_required if no `range` specified_

An array of [predicates](../predicate.md). Once true, the loop will exit.

### `vars`

A [loop variables declaration](../loop-vars.md) that exposes values from the `range` to the call.

## Example

The following op will emit:

```
firstKey firstValue 0
secondKey secondValue 1
```

```yaml
description: An example of a serial loop over a rangable value
run:
  serialLoop:
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

<!--
https://github.com/opctl/opctl/issues/909

The following op will emit:
```
2
1
0
```

```yaml
description: |
  This is an example op.yml file that counts down using a serial loop
inputs:
  counter:
    number:
      default: 2
run:
  serialLoop:
    until:
      - eq: [$(counter), 0]
    run:
      container:
        image: { ref: alpine }
        cmd:
          - sh
          - -c
          - |
            echo $(counter)
            echo "\$(($(counter)-1))" > /file
        files:
          /file: $(counter)
```
-->
