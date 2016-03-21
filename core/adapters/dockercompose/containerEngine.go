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

  containerEngine = containerEngineImpl{
    compositionRoot:compositionRoot,
  }

  return

}

type containerEngineImpl struct {
  compositionRoot compositionRoot
}

func (ce containerEngineImpl) InitDevOp(
devOpName string,
) (err error) {
  return ce.compositionRoot.
  InitDevOpUcExecuter().
  Execute(devOpName)
}

func (ce containerEngineImpl) RunDevOp(
devOpName string,
) (devOpRun models.DevOpRunView, err error) {
  return ce.compositionRoot.
  RunDevOpUcExecuter().
  Execute(devOpName)
}
