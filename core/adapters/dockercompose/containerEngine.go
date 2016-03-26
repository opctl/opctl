package dockercompose

import (
  "github.com/dev-op-spec/engine/core/ports"
  "github.com/dev-op-spec/engine/core/models"
)

func NewContainerEngine(
) (containerEngine ports.ContainerEngine, err error) {

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

func (this _containerEngine) InitDevOp(
devOpName string,
) (err error) {
  return this.compositionRoot.
  InitDevOpUseCase().
  Execute(devOpName)
}

func (this _containerEngine) RunDevOp(
devOpName string,
) (devOpRun models.DevOpRunView, err error) {
  return this.compositionRoot.
  RunDevOpUseCase().
  Execute(devOpName)
}
