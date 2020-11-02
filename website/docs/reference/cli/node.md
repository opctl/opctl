## `opctl node create`
Create an in-process node which inherits current
stderr/stdout/stdin/PGId (process group id) and blocks until killed.

> There can be only one node running at a time on a given machine.

### Options

#### `--data-dir`
Path of dir used to store node data

### Notes

#### lockfile
Upon creation, nodes populate a lockfile at `DATA_DIR/lockfile.pid`
containing their PId (process id).

#### concurrency
Prior to node creation, if a lockfile exists, the existing lock holder
will be liveness tested.

Only in the event the existing lock holder is dead will creation of a
new node occur.

#### debugging
Debugging can be accomplished by running `node create` from a terminal
where it's output is easily monitored.

#### cleanup
During creation, `DATA_DIR` will be
cleaned of any existing events, ops, and temp files/dirs to ensure
the created node starts from a clean slate.

## `opctl node kill`
Kill the running node.