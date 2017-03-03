## `events`

### replaying

Events are persisted to disk and can be replayed (as of v0.1.15).
> events are not held across node restarts; any time a node starts it
> clears its event db.

```shell
$ opctl run some-op

# exit & reopen terminal
$ exit

# previous events still available
$ opctl run events
```

### streaming

Events are streamed in realtime as they occur. They can be
streamed in parallel to any number of terminals.
> behind the scenes, events are delivered over websockets

```shell
# from terminal1
$ opctl events

# from terminal2
$ opctl events

# from terminal3
$ opctl run some-op

# events show up on all terminals simultaneously as they occur
```
