package ports

import (
  "github.com/dev-op-spec/engine/core/models"
)

type ContainerEngine interface {
  InitDevOp(
  pathToDevOpDir string,
  ) (err error)

  RunDevOp(
  pathToDevOpDir string,
  ) (devOpRun models.DevOpRunView, err error)
}
