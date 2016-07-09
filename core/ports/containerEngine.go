package ports

//go:generate counterfeiter -o ../adapters/containerengine/fake/containerEngine.go --fake-name FakeContainerEngine ./ ContainerEngine

import "github.com/opctl/engine/core/logging"

type ContainerEngine interface {
  RunOp(
  correlationId string,
  opArgs map[string]string,
  opBundlePath string,
  opName string,
  opRunId string,
  logger logging.Logger,
  ) (err error)

  KillOpRun(
  correlationId string,
  opBundlePath string,
  opRunId string,
  logger logging.Logger,
  ) (err error)
}
