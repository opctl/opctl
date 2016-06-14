package core

//go:generate counterfeiter -o ./fakeApi.go --fake-name FakeApi ./ Api

import (
  "github.com/opctl/engine/core/models"
  "github.com/opctl/engine/core/ports"
)

type Api interface {
  AddSubOp(
  req models.AddSubOpReq,
  ) (err error)

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
filesys ports.Filesys,
) (api Api, err error) {

  var compositionRoot compositionRoot
  compositionRoot, err = newCompositionRoot(
    containerEngine,
    filesys,
  )
  if (nil != err) {
    return
  }

  api = &_api{
    compositionRoot:compositionRoot,
  }

  return
}

type _api struct {
  compositionRoot compositionRoot
}

func (this _api) AddSubOp(
req models.AddSubOpReq,
) (err error) {
  return this.
  compositionRoot.
  AddSubOpUseCase().
  Execute(req)
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
