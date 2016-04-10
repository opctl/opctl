package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type Api interface {
  AddOp(
  req models.AddOpReq,
  ) (err error)

  AddSubOp(
  req models.AddSubOpReq,
  ) (err error)

  GetEventStream(
  eventChannel chan models.Event,
  ) (err error)

  GetLogForOpRun(
  opRunId string,
  logChannel chan *models.LogEntry,
  ) (err error)

  ListOps(
  projectUrl *models.Url,
  ) (ops []models.OpDetailedView, err error)

  RunOp(
  req models.RunOpReq,
  ) (opRun models.OpRunDetailedView, err error)

  SetDescriptionOfOp(
  req models.SetDescriptionOfOpReq,
  ) (err error)
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

func (this _api) AddOp(
req models.AddOpReq,
) (err error) {
  return this.
  compositionRoot.
  AddOpUseCase().
  Execute(req)
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

func (this _api) GetLogForOpRun(
opRunId string,
logChannel chan *models.LogEntry,
) (err error) {
  return this.
  compositionRoot.
  GetLogForOpRunUseCase().
  Execute(
    opRunId,
    logChannel,
  )
}

func (this _api) ListOps(
projectUrl *models.Url,
) (ops []models.OpDetailedView, err error) {
  return this.
  compositionRoot.
  ListOpsUseCase().
  Execute(projectUrl)
}

func (this _api) RunOp(
req models.RunOpReq,
) (opRun models.OpRunDetailedView, err error) {
  return this.
  compositionRoot.
  RunOpUseCase().
  Execute(
    req,
    make([]*models.Url, 0),
  )
}

func (this _api) SetDescriptionOfOp(
req models.SetDescriptionOfOpReq,
) (err error) {
  return this.
  compositionRoot.
  SetDescriptionOfOpUseCase().
  Execute(
    req,
  )
}
