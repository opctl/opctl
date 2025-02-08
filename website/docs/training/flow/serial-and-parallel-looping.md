---
title: Serial and parallel looping
---

## TLDR;
Opctl supports using [serialLoop](../../reference/opspec/op-directory/op/call/serial-loop.md) and/or [parallelLoop](../../reference/opspec/op-directory/op/call/parallel-loop.md) statements to make parts of your op run in a loop.

## Example
1. Start this op: 
    ```yaml
    name: serialAndParallelLooping
    run:
      serial:
        - parallelLoop:
            range: [1,2,3]
            vars:
              index: $(index)
              value: $(value)
            run:
              container:
                cmd: [echo, "parallelLoop| index: $(index), value: $(value)"]
                image: { ref: ghcr.io/linuxcontainers/alpine }
        - serialLoop:
            range: [1,2,3]
            vars:
              index: $(index)
              value: $(value)
            run:
              container:
                cmd: [echo, "serialLoop| index: $(index), value: $(value)"]
                image: { ref: ghcr.io/linuxcontainers/alpine }
    ```
1. Observe:
   1. for the `parallelLoop` statement, containers for all values in `range` are run in parallel (all at once without order)
   1. for the `serialLoop` statement, containers for all values in `range` are run serially (one by one in order)
