# Implementation

Scenarios are organized hierarchically mimicking the structure of an op as closely as possible to ease
maintenance/discoverability

Each scenario is an op.

Ops containing a scenarios.json file implement one or more test
scenarios (see [scenarios.schema.json](scenarios.schema.json) for schema)
Call scenarios use `scenarios.json`. Add `call.expectOutputs` to check the exact
boxed public output map on the root op's completion event, for example:

```json
"expectOutputs": {"echoed": {"string": "hello"}, "count": {"number": 7}}
```

An empty object asserts that there are no public outputs. Omit the property to
check only the CLI exit status. Output assertions run after an expected successful
CLI call. Use an absolute `/repo/...` local `opRef` so the observer can match its
resolved path and Go ops can find their parent module.

The Go parent-module fixture uses `call.cacheProof: true` to run twice on one
node. Its driver asserts one cold build, zero warm builds, stable executable
bytes/mtime, and fresh plans with changed public inputs, and prints timings.
