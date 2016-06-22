package tcp

/* resources */
const (
  // resource: event-stream
  eventStreamRelUrlTemplate string = "/event-stream"

  // resource: a single liveness
  livenessRelUrlTemplate string = "/liveness"

  // resource: all op-runs
  opRunKillsRelUrlTemplate string = "/op-run-kills"

  // resource: all op-runs
  opRunsRelUrlTemplate string = "/op-runs"
)

/* use cases */
const (
  getLivenessRelUrlTemplate string = livenessRelUrlTemplate
  getEventStreamRelUrlTemplate string = eventStreamRelUrlTemplate
  killOpRunRelUrlTemplate string = opRunKillsRelUrlTemplate
  runOpRelUrlTemplate string = opRunsRelUrlTemplate
)
