---
sidebar_label: create
title: opctl node create
---

```sh
opctl node create
```

Create an in-process node which inherits current
stderr/stdout/stdin/PGId (process group id) and blocks until killed.

> There can be only one node running at a time on a given machine.

## Global Options
see [global options](../global-options.md)

## Notes

### lockfile
Upon creation, nodes populate a lockfile at `DATA_DIR/lockfile.pid`
containing their PId (process id).

### concurrency
Prior to node creation, if a lockfile exists, the existing lock holder
will be liveness tested.

Only in the event the existing lock holder is dead will creation of a
new node occur.
