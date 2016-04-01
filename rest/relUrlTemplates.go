package rest

/* resources */
const (
// resource: a single project
  projectRelUrlTemplate string = "/project/{projectUrl}"

// resource: all project dev ops
  projectDevOpsRelUrlTemplate string = projectRelUrlTemplate + "/dev-ops"

// resource: a single project dev op
  projectDevOpRelUrlTemplate string = projectDevOpsRelUrlTemplate + "/{devOpName}"

// resource: description of a project dev op
  projectDevOpDescriptionRelUrlTemplate string = projectDevOpRelUrlTemplate + "/description"

// resource: all runs of a project dev op
  projectDevOpRunsRelUrlTemplate string = projectDevOpRelUrlTemplate + "/runs"

// resource: all project pipelines
  projectPipelinesRelUrlTemplate string = projectRelUrlTemplate + "/pipelines"

// resource: a single project pipeline
  projectPipelineRelUrlTemplate string = projectPipelinesRelUrlTemplate + "/{pipelineName}"

// resource: description of a project pipeline
  projectPipelineDescriptionRelUrlTemplate string = projectPipelineRelUrlTemplate + "/description"

// resource: all runs of a project pipeline
  projectPipelineRunsRelUrlTemplate string = projectPipelineRelUrlTemplate + "/runs"

// resource: all stages of a project pipeline
  projectPipelineStagesRelUrlTemplate string = projectPipelineRelUrlTemplate + "/stages"
)

/* use cases */
const (
  addDevOpRelUrlTemplate string = projectDevOpsRelUrlTemplate
  addPipelineRelUrlTemplate string = projectPipelinesRelUrlTemplate
  addStageToPipelineRelUrlTemplate string = projectPipelineStagesRelUrlTemplate
  listDevOpsRelUrlTemplate string = projectDevOpsRelUrlTemplate
  listPipelinesRelUrlTemplate string = projectPipelinesRelUrlTemplate
  runDevOpRelUrlTemplate string = projectDevOpRunsRelUrlTemplate
  runPipelineRelUrlTemplate string = projectPipelineRunsRelUrlTemplate
  setDescriptionOfDevOpRelUrlTemplate string = projectDevOpDescriptionRelUrlTemplate
  setDescriptionOfPipelineRelUrlTemplate string = projectPipelineDescriptionRelUrlTemplate
)
