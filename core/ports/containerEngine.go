package ports

import "github.com/dev-op-spec/engine/core/models"

type ContainerEngine interface {
  InitOp(
  pathToOpDir string,
  name string,
  ) (err error)

  RunOp(
  pathToOpDir string,
  name string,
  ) (exitCode int, logChannel chan *models.LogEntry, err error)
}
