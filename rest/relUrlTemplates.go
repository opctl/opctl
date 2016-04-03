package rest

/* resources */
const (
// resource: a single project
  projectRelUrlTemplate string = "/project/{projectUrl}"

// resource: all project operations
  projectOperationsRelUrlTemplate string = projectRelUrlTemplate + "/operations"

// resource: a single project operation
  projectOperationRelUrlTemplate string = projectOperationsRelUrlTemplate + "/{operationName}"

// resource: description of a project operation
  projectOperationDescriptionRelUrlTemplate string = projectOperationRelUrlTemplate + "/description"

// resource: all runs of a project operation
  projectOperationRunsRelUrlTemplate string = projectOperationRelUrlTemplate + "/runs"

// resource: all subOperations of a project operation
  projectOperationSubOperationsRelUrlTemplate string = projectOperationRelUrlTemplate + "/sub-operations"
)

/* use cases */
const (
  addOperationRelUrlTemplate string = projectOperationsRelUrlTemplate
  addSubOperationRelUrlTemplate string = projectOperationSubOperationsRelUrlTemplate
  listOperationsRelUrlTemplate string = projectOperationsRelUrlTemplate
  runOperationRelUrlTemplate string = projectOperationRunsRelUrlTemplate
  setDescriptionOfOperationRelUrlTemplate string = projectOperationDescriptionRelUrlTemplate
)
