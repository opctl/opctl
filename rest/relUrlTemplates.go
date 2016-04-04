package rest

/* resources */
const (
// resource: a single project
  projectRelUrlTemplate string = "/projects/{projectUrl}"

// resource: all project operations
  projectOperationsRelUrlTemplate string = projectRelUrlTemplate + "/operations"

// resource: a single project operation
  projectOperationRelUrlTemplate string = projectOperationsRelUrlTemplate + "/{operationName}"

// resource: description of a project operation
  projectOperationDescriptionRelUrlTemplate string = projectOperationRelUrlTemplate + "/description"

// resource: all subOperations of a project operation
  projectOperationSubOperationsRelUrlTemplate string = projectOperationRelUrlTemplate + "/sub-operations"

// resource: all operation-runs
  operationRunsRelUrlTemplate string = "/operation-runs"
)

/* use cases */
const (
  addOperationRelUrlTemplate string = projectOperationsRelUrlTemplate
  addSubOperationRelUrlTemplate string = projectOperationSubOperationsRelUrlTemplate
  listOperationsRelUrlTemplate string = projectOperationsRelUrlTemplate
  runOperationRelUrlTemplate string = operationRunsRelUrlTemplate
  setDescriptionOfOperationRelUrlTemplate string = projectOperationDescriptionRelUrlTemplate
)
