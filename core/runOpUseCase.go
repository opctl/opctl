package core

//go:generate counterfeiter -o ./fakeRunOpUseCase.go --fake-name fakeRunOpUseCase ./ runOpUseCase

import (
  "github.com/opctl/engine/core/models"
  "github.com/opspec-io/sdk-golang"
  "path/filepath"
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
opspecSdk opspec.Sdk,
uniqueStringFactory uniqueStringFactory,
) runOpUseCase {

  return &_runOpUseCase{
    opRunner:opRunner,
    opspecSdk:opspecSdk,
    uniqueStringFactory:uniqueStringFactory,
  }

}

type _runOpUseCase struct {
  opRunner            opRunner
  opspecSdk           opspec.Sdk
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

  opCollection, err := this.opspecSdk.GetCollection(
    filepath.Dir(req.OpUrl),
  )
  if (nil != err) {
    return
  }

  opRunId, err = this.opRunner.Run(
    correlationId,
    req.Args,
    opCollection.Name,
    req.OpUrl,
    "",
  )

  return

}
