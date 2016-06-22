package dockercompose

import (
  "github.com/opctl/engine/core/ports"
  "github.com/opctl/engine/core/logging"
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

func (this _containerEngine) InitOp(
opBundlePath string,
opName string,
) (err error) {

  return this.compositionRoot.
  InitOpUseCase().
  Execute(
    opBundlePath,
    opName,
  )

}

func (this _containerEngine) RunOp(
opArgs map[string]string,
correlationId string,
opBundlePath string,
opName string,
opNamespace string,
logger logging.Logger,
) (
exitCode int,
err error,
) {

  return this.compositionRoot.
  RunOpUseCase().
  Execute(
    opArgs,
    correlationId,
    opBundlePath,
    opName,
    opNamespace,
    logger,
  )

}

func (this _containerEngine) KillOpRun(
correlationId string,
opBundlePath string,
opNamespace string,
logger logging.Logger,
) (err error) {

  return this.compositionRoot.
  KillOpRunUseCase().
  Execute(
    correlationId,
    opBundlePath,
    opNamespace,
    logger,
  )

}

