## `node create` (since v0.1.15)

Create an in-process node which inherits current
stderr/stdout/stdin/PGId (process group id) and blocks until killed.

> There can be only one node running at a time on a given machine.


## Notes

### lockfile

Upon creation, nodes populate a lockfile at
[per user app data](https://github.com/appdataspec/spec/blob/master/index.md#per-user-app-data)/lockfile.pid
containing their PId (process id).

### concurrency

Prior to node creation, if a lockfile exists, the existing lock holder
will be liveness tested.

Only in the event the existing lock holder is dead will creation of a
new node occur.

### debugging

Debugging can be accomplished by running `node create` from a terminal
where it's output is easily monitored.

### cleanup

During creation,
[node root path](../../node/filesystem.md#node-root-path) will be
cleaned of any existing events, packages, and temp files/dirs to ensure
the created node starts from a clean slate.
