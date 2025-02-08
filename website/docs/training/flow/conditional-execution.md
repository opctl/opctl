---
title: Conditional execution
---

## TLDR;
Opctl supports using [if](../../reference/opspec/op-directory/op/call/index.md#if) statements and [predicates](../../reference/opspec/op-directory/op/call/predicate.md) to make parts of your op run conditionally.

## Example
1. Start this op: 
    ```yaml
    name: conditionalExecution
    inputs:
      shouldRunContainer:
        description: whether to run the container or not
        boolean: {}
    run:
      if:
        - eq: [true, $(shouldRunContainer)]
      container:
        cmd: [echo, 'hello!']
        image: { ref: ghcr.io/linuxcontainers/alpine }
    ```
1. When prompted, enter `true` or `false`
1. Observe you only see the container run and `hello!` logged when you enter `true`.
