package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type Api interface {
  AddDevOp(
  req models.AddDevOpReq,
  ) (err error)

  AddPipeline(
  req models.AddPipelineReq,
  ) (err error)

  AddStageToPipeline(
  req models.AddStageToPipelineReq,
  ) (err error)

  ListDevOps(
  projectUrl *models.ProjectUrl,
  ) (devOps []models.DevOpView, err error)

  ListPipelines(
  projectUrl *models.ProjectUrl,
  ) (pipelines []models.PipelineView, err error)

  RunDevOp(
  req models.RunDevOpReq,
  ) (devOpRun models.DevOpRunView, err error)

  RunPipeline(
  req models.RunPipelineReq,
  ) (pipelineRun models.PipelineRunView, err error)

  SetDescriptionOfDevOp(
  req models.SetDescriptionOfDevOpReq,
  ) (err error)

  SetDescriptionOfPipeline(
  req models.SetDescriptionOfPipelineReq,
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

func (this _api) AddDevOp(
req models.AddDevOpReq,
) (err error) {
  return this.
  compositionRoot.
  AddDevOpUseCase().
  Execute(req)
}

func (this _api) AddPipeline(
req models.AddPipelineReq,
) (err error) {
  return this.
  compositionRoot.
  AddPipelineUseCase().
  Execute(req)
}

func (this _api) AddStageToPipeline(
req models.AddStageToPipelineReq,
) (err error) {
  return this.
  compositionRoot.
  AddStageToPipelineUseCase().
  Execute(req)
}

func (this _api) ListDevOps(
projectUrl *models.ProjectUrl,
) (devOps []models.DevOpView, err error) {
  return this.
  compositionRoot.
  ListDevOpsUseCase().
  Execute(projectUrl)
}

func (this _api) ListPipelines(
projectUrl *models.ProjectUrl,
) (pipelines []models.PipelineView, err error) {
  return this.
  compositionRoot.
  ListPipelinesUseCase().
  Execute(projectUrl)
}

func (this _api) RunDevOp(
req models.RunDevOpReq,
) (devOpRun models.DevOpRunView, err error) {
  return this.
  compositionRoot.
  RunDevOpUseCase().
  Execute(req)
}

func (this _api) RunPipeline(
req models.RunPipelineReq,
) (pipelineRun models.PipelineRunView, err error) {
  return this.
  compositionRoot.
  RunPipelineUseCase().
  Execute(
    req,
    make([]string, 0),
  )
}

func (this _api) SetDescriptionOfDevOp(
req models.SetDescriptionOfDevOpReq,
) (err error) {
  return this.
  compositionRoot.
  SetDescriptionOfDevOpUseCase().
  Execute(
    req,
  )
}

func (this _api) SetDescriptionOfPipeline(
req models.SetDescriptionOfPipelineReq,
) (err error) {
  return this.
  compositionRoot.
  SetDescriptionOfPipelineUseCase().
  Execute(
    req,
  )
}
