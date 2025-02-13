---
title: Serial and parallel execution
---

## TLDR;
Opctl supports using [serial](../../reference/opspec/op-directory/op/call/index.md#serial) and/or [parallel](../../reference/opspec/op-directory/op/call/index.md#parallel) statements to make parts of your op run serially (one by one in order) or in parallel (all at once without order).

## Example
1. Start this op: 
    ```yaml
    name: serialAndParallelExecution
    run:
      serial:
        - parallel:
            - container:
                cmd: [echo, "parallel[0]"]
                image: { ref: ghcr.io/linuxcontainers/alpine }
            - container:
                cmd: [echo, "parallel[1]"]
                image: { ref: ghcr.io/linuxcontainers/alpine }
        - serial:
            - container:
                cmd: [echo, "serial[0]"]
                image: { ref: ghcr.io/linuxcontainers/alpine }
            - container:
                cmd: [echo, "serial[1]"]
                image: { ref: ghcr.io/linuxcontainers/alpine }
    ```
2. Observe:
   1. for the `parallel` statement, containers are run in parallel (all at once without order)
   2. for the `serial` statement, containers are run serially (one by one in order)
