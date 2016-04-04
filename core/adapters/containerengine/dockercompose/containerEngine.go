package dockercompose

import (
  "github.com/dev-op-spec/engine/core/ports"
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

func (this _containerEngine) InitOperation(
pathToOperationDir string,
) (err error) {
  return this.compositionRoot.
  InitOperationUseCase().
  Execute(pathToOperationDir)
}

func (this _containerEngine) RunOperation(
pathToOperationDir string,
) (exitCode int, err error) {
  return this.compositionRoot.
  RunOperationUseCase().
  Execute(pathToOperationDir)
}
