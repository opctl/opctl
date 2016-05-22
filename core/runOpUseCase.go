package core

//go:generate counterfeiter -o ./fakeRunOpUseCase.go --fake-name fakeRunOpUseCase ./ runOpUseCase

import (
  "github.com/opctl/engine/core/models"
)

type runOpUseCase interface {
  Execute(
  req models.RunOpReq,
  ) (
  opRunId string,
  correlationId string,
  err error,
  )
}

func newRunOpUseCase(
opRunner opRunner,
uniqueStringFactory uniqueStringFactory,
) runOpUseCase {

  return &_runOpUseCase{
    opRunner:opRunner,
    uniqueStringFactory:uniqueStringFactory,
  }

}

type _runOpUseCase struct {
  opRunner            opRunner
  uniqueStringFactory uniqueStringFactory
}

func (this _runOpUseCase) Execute(
req models.RunOpReq,
) (
opRunId string,
correlationId string,
err error,
) {

  correlationId, err = this.uniqueStringFactory.Construct()
  if (nil != err) {
    return
  }

  opRunId, err = this.opRunner.Run(
    correlationId,
    req.OpUrl,
    "",
  )

  return

}
