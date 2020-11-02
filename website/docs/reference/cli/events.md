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
   opctl events
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