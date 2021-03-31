---
sidebar_label: Overview
title: Call
---

A call is an object that defines a single call within an ops call graph. Opctl supports several types of calls and several advanced attributes to add logic into your op.

The leaves of a call graph are [containers](container/index.md) that run programs to do work. Each node in the call graph will end once all child containers exit successfully or a single container exits with a failure, which will cause still running containers to be killed.

Containers and call graph nodes communicate with each other by producing and emitting [data](../index.md) or by reading and writing shared [files](../types/file.md) and [directories](../types/dir.md) (which are passed by reference).

## Basic propreties

### `description`

A human friendly description of the parameter, written as a [markdown string](../markdown.md).

### `name`

The name is an [identifier](../identifier.md) used to identify the call in a UI and used for [`needs`](#needs) logic in sibling calls.

## Types of calls properties

Every object in your call graph must declare one of the following properties to define how and what it runs. These can directly run an op or container or can sequence further calls in serial, parallel, or by looping.

### `container`

A [container call](container/index.md) runs a container.

### `op`

An [op call](container/index.md) runs an external op.

### `parallel`

A parallel call is an array of [calls](index.md) that are executed concurrently in parallel, with no defined order.

### `parallelLoop`

A [parallel loop call](parallel-loop) defines a call executed multiple times in parallel.

### `serial`

A serial call is an array of [calls](index.md) that are executed serially, one after another.

### `serialLoop`

A [serial loop call](serial-loop) defines a call executed multiple times repeatedly.

## Logical properties

Calls can optionally define additional properties to introduce conditional logic into their execution.

### `if`

An if property on a call is an array of [predicates](../predicate.md). If all predicates are true, the call is executed, otherwise it is skipped.

```yaml
if:
  - eq: [true, $(test)]
  - exists: $(value)
```

### `needs`

`needs` allows introducing dependencies between [parallel](#parallel) calls. The value of this property is an array of sibling [names](#name). When all calls that need a given name end, the named call will be killed.

Needs cannot be used in parents, children, or cousin calls (they must be within the same `parallel` call). Opctl will ignore needs that don't satisfy this limitation.

```yaml
description: |
  `systemUnderTest` will be shutdown after 1 second because it's no longer needed by the second container.
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
