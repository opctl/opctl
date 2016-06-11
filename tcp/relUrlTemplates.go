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

  // resource: a single project
  projectRelUrlTemplate string = "/projects/{projectUrl}"

  // resource: all project ops
  projectOpsRelUrlTemplate string = projectRelUrlTemplate + "/ops"

  // resource: a single project op
  projectOpRelUrlTemplate string = projectOpsRelUrlTemplate + "/{opName}"

  // resource: all subOps of a project op
  projectOpSubOpsRelUrlTemplate string = projectOpRelUrlTemplate + "/sub-ops"
)

/* use cases */
const (
  addOpRelUrlTemplate string = projectOpsRelUrlTemplate
  addSubOpRelUrlTemplate string = projectOpSubOpsRelUrlTemplate
  getLivenessRelUrlTemplate string = livenessRelUrlTemplate
  getEventStreamRelUrlTemplate string = eventStreamRelUrlTemplate
  killOpRunRelUrlTemplate string = opRunKillsRelUrlTemplate
  listOpsRelUrlTemplate string = projectOpsRelUrlTemplate
  runOpRelUrlTemplate string = opRunsRelUrlTemplate
)
