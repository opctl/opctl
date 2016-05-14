package ports

//go:generate counterfeiter -o ../adapters/containerengine/fake/containerEngine.go --fake-name FakeContainerEngine ./ ContainerEngine

import "github.com/dev-op-spec/engine/core/logging"

type ContainerEngine interface {
  InitOp(
  pathToOpDir string,
  name string,
  ) (err error)

  RunOp(
  correlationId string,
  pathToOpDir string,
  name string,
  logger logging.Logger,
  ) (exitCode int, err error)

  KillOpRun(
  correlationId string,
  pathToOpDir string,
  logger logging.Logger,
  ) (err error)
}
