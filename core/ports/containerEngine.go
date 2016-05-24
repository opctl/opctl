package ports

//go:generate counterfeiter -o ../adapters/containerengine/fake/containerEngine.go --fake-name FakeContainerEngine ./ ContainerEngine

import "github.com/opctl/engine/core/logging"

type ContainerEngine interface {
  InitOp(
  pathToOpDir string,
  name string,
  ) (err error)

  RunOp(
  args map[string]string,
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
