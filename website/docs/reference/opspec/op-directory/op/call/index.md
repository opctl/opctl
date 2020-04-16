---
sidebar_label: Index
title: Call [object]
---

An object defining a call in an operations call graph.

## Properties
- must have exactly one of
  - [container](#container)
  - [op](#op)
  - [parallel](#parallel)
  - [parallelLoop](#parallelloop)
  - [serial](#serial)
  - [serialLoop](#serialloop)
- may have
  - [if](#if)
  - [name](#name)
  - [needs](#needs)

### container
A [container-call [object]](container/index.md) defining a container to run.

### op
An [op-call [object]](op.md) defining an op to run.

### parallel
An array of [call [object]](index.md)s defining calls run in parallel (all at once without order).

### parallelLoop
A [parallel-loop-call [object]](parallel-loop.md) defining a call loop in which all iterations happen in parallel (all at once without order).

### serial
An array of [call [object]](index.md)s defining calls run in serial (one after another in order).

### serialLoop
A [serial-loop-call [object]](serial-loop.md) defining a call loop in which each iteration happens in serial (one after another in order)

### if
An array of [predicate [object]](predicate.md)s which must all be true for the call to take place.

### name
An [identifier [string]](../identifier.md) used to identify the call in UI's or [needs](#needs) of sibling calls.

### needs
An array of [identifier [string]](../identifier.md)s identifying calls needed by the current call. When the named calls are no longer needed (by this or any other call), they will be killed.

> note: needed calls and the current call MUST be children of the same parallel block. If not, the need will be ignored.

#### Example Needs (Integration Test)
```yaml
name: integration-test
description: systemUnderTest will be shutdown after 1 second because it's no longer needed.
run:
  parallel:
    - name: systemUnderTest
      container:
        image: {ref: alpine}
        cmd: [sleep, 100000]
    - container:
        image: {ref: alpine}
        cmd: [sleep, 1]
      needs:
        - systemUnderTest
```