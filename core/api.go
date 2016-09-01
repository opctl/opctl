package core

//go:generate counterfeiter -o ./fakeApi.go --fake-name FakeApi ./ Api

import (
  "github.com/opspec-io/engine/core/models"
  "github.com/opspec-io/engine/core/ports"
)

type Api interface {
  GetEventStream(
  eventChannel chan models.Event,
  ) (err error)

  KillOpRun(
  req models.KillOpRunReq,
  ) (
  correlationId string,
  err error,
  )

  RunOp(
  req models.RunOpReq,
  ) (
  opRunId string,
  correlationId string,
  err error,
  )
}

func New(
containerEngine ports.ContainerEngine,
) (api Api) {

  api = &_api{
    compositionRoot:
    newCompositionRoot(
      containerEngine,
    ),
  }

  return
}

type _api struct {
  compositionRoot compositionRoot
}

func (this _api) GetEventStream(
eventChannel chan models.Event,
) (err error) {
  return this.
  compositionRoot.
  GetEventStreamUseCase().
  Execute(eventChannel)
}

func (this _api) KillOpRun(
req models.KillOpRunReq,
) (
correlationId string,
err error,
) {
  return this.
  compositionRoot.
  KillOpRunUseCase().
  Execute(req)
}

func (this _api) RunOp(
req models.RunOpReq,
) (
opRunId string,
correlationId string,
err error,
) {
  return this.
  compositionRoot.
  RunOpUseCase().
  Execute(req)
}
