package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type Api interface {
  AddOperation(
  req models.AddOperationReq,
  ) (err error)

  AddSubOperation(
  req models.AddSubOperationReq,
  ) (err error)

  ListOperations(
  projectUrl *models.ProjectUrl,
  ) (operations []models.OperationDetailedView, err error)

  RunOperation(
  req models.RunOperationReq,
  ) (operationRun models.OperationRunDetailedView, err error)

  SetDescriptionOfOperation(
  req models.SetDescriptionOfOperationReq,
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

func (this _api) AddOperation(
req models.AddOperationReq,
) (err error) {
  return this.
  compositionRoot.
  AddOperationUseCase().
  Execute(req)
}

func (this _api) AddSubOperation(
req models.AddSubOperationReq,
) (err error) {
  return this.
  compositionRoot.
  AddSubOperationUseCase().
  Execute(req)
}

func (this _api) ListOperations(
projectUrl *models.ProjectUrl,
) (operations []models.OperationDetailedView, err error) {
  return this.
  compositionRoot.
  ListOperationsUseCase().
  Execute(projectUrl)
}

func (this _api) RunOperation(
req models.RunOperationReq,
) (operationRun models.OperationRunDetailedView, err error) {
  return this.
  compositionRoot.
  RunOperationUseCase().
  Execute(
    req,
    make([]string, 0),
  )
}

func (this _api) SetDescriptionOfOperation(
req models.SetDescriptionOfOperationReq,
) (err error) {
  return this.
  compositionRoot.
  SetDescriptionOfOperationUseCase().
  Execute(
    req,
  )
}
