## `events`

listen to node events.

> if a node isn't running, one will be automatically created.

## Notes

### replaying

Events are persisted to disk and can be replayed (since v0.1.15).
> events are not held across node restarts; any time a node starts it
> clears its event db.

example:
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

example:
```shell
# from terminal1
$ opctl events

# from terminal2
$ opctl events

# from terminal3
$ opctl run some-op

# events show up on all terminals simultaneously as they occur
```
