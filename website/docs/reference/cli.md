---
title: Commands
sidebar_label: CLI
---

## `opctl`

### Options

#### `--no-color`

To disable color, include a `--no-color` flag w/ your
command.
> this may increase readability in environments not supporting
> color escape codes or piping output to another program

##### Examples

```sh
opctl --no-color events
```

#### `-h` or `--help`

For context specific help, include a `-h` (or `--help`) flag w/ your
command.

##### Examples

```sh
opctl node create -h

Usage: opctl node create

Creates a node
```

## `opctl events`

listen to node events.

> if a node isn't running, one will be automatically created.


### Examples

#### Event Replay
Events are persisted to disk and can be replayed.
> events are not held across node restarts; any time a node starts it
> clears its event db.

1. open terminal & run an op so we have some events
   ```sh
   opctl run github.com/opspec-pkgs/uuid.v4.generate#1.1.0
   ```

1. exit terminal
   ```sh
   exit
   ```

1. re open terminal & replay events
   ```sh
   opctl run events
   ```

#### Event Streaming
Events are streamed in realtime as they occur. They can be streamed in parallel to any number of terminals.
> behind the scenes, events are delivered over websockets

1. open multiple terminals & open event stream on each
   ```sh
   opctl events
   ```

1. open another terminal & run an op; watch events show up on all terminals simultaneously in real-time
   ```sh
   opctl run some-op
   ```

## `opctl ls [DIR_REF]`

List ops in a directory.

### Arguments

#### `DIR_REF` *default: `.opspec`*

Reference to dir ops will be listed from (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)

### Examples

#### `.opspec` dir

lists ops from `./.opspec`

```sh
opctl ls
```

#### remote dir

lists ops from [github.com/opctl/opctl#0.1.24](https://github.com/opctl/opctl/tree/0.1.24)

```sh
opctl ls github.com/opctl/opctl#0.1.24/
```

## `opctl run [OPTIONS] OP_REF`

Start and wait on an op

> if a node isn't running, one will be automatically created

### Arguments

#### `OP_REF`

Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`)

### Options

#### `-a`

Explicitly pass args to op in format `-a NAME1=VALUE1 -a NAME2=VALUE2`

#### `--arg-file` *default: `.opspec/args.yml`*

Read in a file of args in yml format

### Examples

#### local op ref w/out args

```sh
opctl run myop
```

#### remote op ref w/ args

```sh
opctl run -a apiToken="my-token" -a channelName="my-channel" -a msg="hello!" github.com/opspec-pkgs/slack.chat.post-message#0.1.1
```

### Notes

#### op source username/password prompt

If auth w/ the op source fails the cli will (re)prompt for username &
password.

> in non-interactive terminals, the cli will note that it can't prompt
> due to being in a non-interactive terminal and exit with a non zero
> exit code.

#### input sources

Input sources are checked according to the following precedence:

- arg provided via `-a` option
- arg file
- env var
- default
- prompt

#### input prompts

Inputs which are invalid or missing will result in the cli prompting for
them.

> in non-interactive terminals, the cli will provide details about the
> invalid or missing input, note that it's giving up due to being in a
> non-interactive terminal and exit with a non zero exit code.

example:

```sh

-
  Please provide value for parameter.
  Name: version
  Description: version of app being compiled
-
```

##### validation

When inputs don't meet constraints, the cli will (re)prompt for the
input until a satisfactory value is obtained.

#### containers

##### image

###### image layer caching

All pulled image layers will be cached

###### image updates

Prior to container creation, updates to the referenced image will be
pulled and applied.

If checking for or applying updated image layers fails, graceful
fallback to cached image layers will occur

##### networking

All containers created by opctl will be attached to a single managed
network.

> the network is visible from `docker network ls` as `opctl`.

##### cleanup

Containers will be removed as they exit.

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

## `opctl op create [OPTIONS] NAME`

Creates an op

### Arguments

#### `NAME`
Name of the op

### Options

#### `-d` or `--description`
Description of the op

#### `--path` *default: `.opspec`*
Path to create the op at

### Examples

```sh
opctl op create -d "my awesome op description" --path some/path my-awesome-op-name
```

## `op install [OPTIONS] OP_REF`

Installs an op

### Arguments

#### `OP_REF`

Op reference (`host/path/repo#tag`, or `host/path/repo#tag/path`)

### Options

#### `--path` *default: `.opspec/OP_REF`*

Path to install the op at

#### `-u` or `--username`

Username used to auth w/ the op source

#### `-p` or `--password`

Password used to auth w/ the op source

### Examples

```sh
opctl op install -u someUser -p somePass host/path/repo#tag
```

### Notes

#### op source username/password prompt

If auth w/ the op source fails the cli will (re)prompt for username &
password.

> in non-interactive terminals, the cli will note that it can't prompt
> and exit with a non zero exit code.

## `op kill OP_ID`

Kill an op. 

### Arguments

#### `OP_ID`
Id of the op to kill

## `op validate OP_REF`

Validates an op according to:

- existence of `op.yml`
- validity of `op.yml` (per
  [schema](https://opctl.io/0.1.6/op.yml.schema.json))

### Arguments

#### `OP_REF`

Op reference (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`).

### Examples

```sh
opctl op validate myop
```

### Notes

#### op source username/password prompt

If auth w/ the op source fails the cli will (re)prompt for username & password.

> in non-interactive terminals, the cli will note that it can't prompt and exit with a non zero exit code.

## `opctl self-update [OPTIONS]`

Updates the current version of opctl.

> if a node is running, it will be automatically killed

### Options

##### `-c` or `--channel` *default: `stable`*
The release channel to update from

- `stable`
- `beta` (smoke tested alpha channel)
- `alpha` (all bets are off)

### Examples

get latest stable release
```sh
opctl self-update
# output: Updated to new version: 0.1.24!
```

play around w/ latest beta release
 ```sh
opctl self-update -c beta
# output: Updated to new version: 0.1.24-beta.1!
```

play times over; switch back to latest stable release
```sh
opctl self-update
# output: Updated to new version: 0.1.24!
```

## `ui [MOUNT_REF]`

Opens the opctl web UI to the current working directory.

### Arguments

#### `MOUNT_REF` *default: `.`*
Reference to mount (either `relative/path`, `/absolute/path`, `host/path/repo#tag`, or `host/path/repo#tag/path`).

### Examples

Open web UI to current working directory
```sh
opctl ui
```

Open web UI to remote op (github.com/opspec-pkgs/_.op.create#3.3.1)
```sh
opctl ui github.com/opspec-pkgs/_.op.create#3.3.1
```