package dockercompose

import (
  "github.com/opspec-io/engine/core/ports"
  "github.com/opspec-io/engine/core"
)

func New(
) (
containerEngine ports.ContainerEngine,
err error,
) {

  var compositionRoot compositionRoot
  compositionRoot, err = newCompositionRoot()
  if (nil != err) {
    return
  }

  containerEngine = _containerEngine{
    compositionRoot:compositionRoot,
  }

  return

}

type _containerEngine struct {
  compositionRoot compositionRoot
}

func (this _containerEngine) StartContainer(
opRunArgs map[string]string,
opBundlePath string,
opName string,
opRunId string,
eventPublisher core.EventPublisher,
rootOpRunId string,
) (err error) {

  return this.compositionRoot.
    StartContainerUseCase().
    Execute(
    opRunArgs,
    opBundlePath,
    opName,
    opRunId,
    eventPublisher,
    rootOpRunId,
  )

}

func (this _containerEngine) EnsureContainerKilled(
opBundlePath string,
opRunId string,
) {

  this.compositionRoot.
    EnsureContainerRemovedUseCase().
    Execute(
    opBundlePath,
    opRunId,
  )

}

