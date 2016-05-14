package dockercompose

import (
  "github.com/dev-op-spec/engine/core/ports"
  "github.com/dev-op-spec/engine/core/logging"
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

func (this _containerEngine) InitOp(
pathToOpDir string,
opName string,
) (err error) {

  return this.compositionRoot.
  InitOpUseCase().
  Execute(
    pathToOpDir,
    opName,
  )

}

func (this _containerEngine) RunOp(
correlationId string,
pathToOpDir string,
opName string,
logger logging.Logger,
) (exitCode int, err error) {

  return this.compositionRoot.
  RunOpUseCase().
  Execute(
    correlationId,
    pathToOpDir,
    opName,
    logger,
  )

}

func (this _containerEngine) KillOpRun(
correlationId string,
pathToOpDir string,
logger logging.Logger,
) (err error) {

  return this.compositionRoot.
  KillOpRunUseCase().
  Execute(
    correlationId,
    pathToOpDir,
    logger,
  )

}

