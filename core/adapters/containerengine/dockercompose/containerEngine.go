package dockercompose

import (
  "github.com/opspec-io/engine/core/ports"
  "github.com/opspec-io/engine/core/logging"
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

func (this _containerEngine) RunOp(
correlationId string,
opArgs map[string]string,
opBundlePath string,
opName string,
opRunId string,
logger logging.Logger,
) (err error) {

  return this.compositionRoot.
  RunOpUseCase().
  Execute(
    correlationId,
    opArgs,
    opBundlePath,
    opName,
    opRunId,
    logger,
  )

}

func (this _containerEngine) KillOpRun(
correlationId string,
opBundlePath string,
opRunId string,
logger logging.Logger,
) (err error) {

  return this.compositionRoot.
  KillOpRunUseCase().
  Execute(
    correlationId,
    opBundlePath,
    opRunId,
    logger,
  )

}

