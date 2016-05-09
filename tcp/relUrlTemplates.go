package tcp

/* resources */
const (
  // resource: a single liveness
  livenessRelUrlTemplate string = "/liveness"

  // resource: a single project
  projectRelUrlTemplate string = "/projects/{projectUrl}"

  // resource: all project ops
  projectOpsRelUrlTemplate string = projectRelUrlTemplate + "/ops"

  // resource: a single project op
  projectOpRelUrlTemplate string = projectOpsRelUrlTemplate + "/{opName}"

  // resource: description of a project op
  projectOpDescriptionRelUrlTemplate string = projectOpRelUrlTemplate + "/description"

  // resource: all subOps of a project op
  projectOpSubOpsRelUrlTemplate string = projectOpRelUrlTemplate + "/sub-ops"

  // resource: all op-runs
  opRunsRelUrlTemplate string = "/op-runs"

  // resource: event-stream
  eventStreamRelUrlTemplate string = "/event-stream"
)

/* use cases */
const (
  addOpRelUrlTemplate string = projectOpsRelUrlTemplate
  addSubOpRelUrlTemplate string = projectOpSubOpsRelUrlTemplate
  getLivenessRelUrlTemplate string = livenessRelUrlTemplate
  getEventStreamRelUrlTemplate string = eventStreamRelUrlTemplate
  listOpsRelUrlTemplate string = projectOpsRelUrlTemplate
  runOpRelUrlTemplate string = opRunsRelUrlTemplate
  setDescriptionOfOpRelUrlTemplate string = projectOpDescriptionRelUrlTemplate
)
