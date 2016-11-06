package containerengine

//go:generate counterfeiter -o engines/fake/containerEngine.go --fake-name ContainerEngine ./ ContainerEngine

import (
  "github.com/opspec-io/engine/util/eventing"
)

type ContainerEngine interface {
  StartContainer(
  opRunArgs map[string]string,
  opBundlePath string,
  opName string,
  opRunId string,
  eventPublisher eventing.EventPublisher,
  rootOpRunId string,
  ) (err error)

  EnsureContainerRemoved(
  opBundlePath string,
  opRunId string,
  )
}
