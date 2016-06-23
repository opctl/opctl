package ports

//go:generate counterfeiter -o ../adapters/containerengine/fake/containerEngine.go --fake-name FakeContainerEngine ./ ContainerEngine

import "github.com/opctl/engine/core/logging"

type ContainerEngine interface {
  InitOp(
  opBundlePath string,
  name string,
  ) (err error)

  RunOp(
  correlationId string,
  opArgs map[string]string,
  opBundlePath string,
  opName string,
  opNamespace string,
  logger logging.Logger,
  ) (exitCode int, err error)

  KillOpRun(
  correlationId string,
  opBundlePath string,
  opNamespace string,
  logger logging.Logger,
  ) (err error)
}
