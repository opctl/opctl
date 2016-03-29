package dockercompose

import (
  "github.com/dev-op-spec/engine/core/ports"
  "github.com/dev-op-spec/engine/core/models"
)

func New(
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
pathToDevOpDir string,
) (err error) {
  return this.compositionRoot.
  InitDevOpUseCase().
  Execute(pathToDevOpDir)
}

func (this _containerEngine) RunDevOp(
pathToDevOpDir string,
) (devOpRun models.DevOpRunView, err error) {
  return this.compositionRoot.
  RunDevOpUseCase().
  Execute(pathToDevOpDir)
}
