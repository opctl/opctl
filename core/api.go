package core

//go:generate counterfeiter -o ./fakeApi.go --fake-name FakeApi ./ Api

import (
  "github.com/opspec-io/sdk-golang/models"
)

type Api interface {
  GetEventStream(
  eventChannel chan models.Event,
  ) (err error)

  KillOpRun(
  req models.KillOpRunReq,
  ) (
  err error,
  )

  StartOpRun(
  req models.StartOpRunReq,
  ) (
  opRunId string,
  err error,
  )
}

func New(
containerEngine ContainerEngine,
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
err error,
) {
  return this.
  compositionRoot.
    KillOpRunUseCase().
    Execute(req)
}

func (this _api) StartOpRun(
req models.StartOpRunReq,
) (
opRunId string,
err error,
) {
  return this.
  compositionRoot.
    StartOpRunUseCase().
    Execute(req)
}
