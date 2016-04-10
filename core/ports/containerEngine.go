package ports

import "github.com/dev-op-spec/engine/core/logging"

type ContainerEngine interface {
  InitOp(
  pathToOpDir string,
  name string,
  ) (err error)

  RunOp(
  pathToOpDir string,
  name string,
  logger logging.Logger,
  ) (exitCode int, err error)
}
