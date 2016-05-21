package core

//go:generate counterfeiter -o ./fakeKillOpRunUseCase.go --fake-name fakeKillOpRunUseCase ./ killOpRunUseCase

import (
  "github.com/open-devops/engine/core/models"
)

type killOpRunUseCase interface {
  Execute(
  req models.KillOpRunReq,
  ) (
  correlationId string,
  err error,
  )
}

func newKillOpRunUseCase(
opRunner opRunner,
uniqueStringFactory uniqueStringFactory,
) killOpRunUseCase {

  return &_killOpRunUseCase{
    opRunner:opRunner,
    uniqueStringFactory:uniqueStringFactory,
  }

}

type _killOpRunUseCase struct {
  opRunner            opRunner
  uniqueStringFactory uniqueStringFactory
}

func (this _killOpRunUseCase) Execute(
req models.KillOpRunReq,
) (
correlationId string,
err error,
) {

  correlationId, err = this.uniqueStringFactory.Construct()
  if (nil != err) {
    return
  }

  err = this.opRunner.Kill(
    correlationId,
    req.OpRunId,
  )

  return

}
