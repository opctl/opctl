`util` contains modules which aren't specific to the SDK & "could" be
used generically.

Placing them here keeps the SDK specific codebase minimal.

Util modules should be developed w/ the mindset they could be published
if it ever made sense to which means:

- they include a clear README.md w/ problem statement
- they have a minimal & clean API
