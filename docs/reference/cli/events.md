## `events`

listen to node events.

> if a node isn't running, one will be automatically created.

## Notes

### replaying

Events are persisted to disk and can be replayed (since v0.1.15).
> events are not held across node restarts; any time a node starts it
> clears its event db.

example:

step 1: open terminal & generate some events by running an op

```shell
opctl run some-op
```

step 2: exit terminal

```shell
exit
```

step 3: re open terminal & replay events

```shell
opctl run events
```

### streaming

Events are streamed in realtime as they occur. They can be streamed in
parallel to any number of terminals.
> behind the scenes, events are delivered over websockets

example:

step 1: open multiple terminals & open event stream on each

```shell
opctl events
```

step 2: open another terminal & run an op; watch events show up on all
terminals simultaneously in real-time

```shell
opctl run some-op
```

