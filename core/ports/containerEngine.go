package ports

import (
  "github.com/dev-op-spec/engine/core/models"
)

type ContainerEngine interface {
  InitOperation(
  pathToOperationDir string,
  ) (err error)

  RunOperation(
  pathToOperationDir string,
  ) (operationRun models.OperationRunDetailedView, err error)
}
