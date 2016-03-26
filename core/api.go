package core

import "github.com/dev-op-spec/engine/core/models"

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
  ) (devOps []models.DevOpView, err error)

  ListPipelines(
  ) (pipelines []models.PipelineView, err error)

  RunDevOp(
  devOpName string,
  ) (devOpRun models.DevOpRunView, err error)

  RunPipeline(
  pipelineName string,
  ) (pipelineRun models.PipelineRunView, err error)

  SetDescriptionOfDevOp(
  req models.SetDescriptionOfDevOpReq,
  ) (err error)

  SetDescriptionOfPipeline(
  req models.SetDescriptionOfPipelineReq,
  ) (err error)
}

func New(
) (api Api, err error) {

  var compositionRoot compositionRoot
  compositionRoot, err = newCompositionRoot()
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
) (devOps []models.DevOpView, err error) {
  return this.
  compositionRoot.
  ListDevOpsUseCase().
  Execute()
}

func (this _api) ListPipelines(
) (pipelines []models.PipelineView, err error) {
  return this.
  compositionRoot.
  ListPipelinesUseCase().
  Execute()
}

func (this _api) RunDevOp(
devOpName string,
) (devOpRun models.DevOpRunView, err error) {
  return this.
  compositionRoot.
  RunDevOpUseCase().
  Execute(devOpName)
}

func (this _api) RunPipeline(
pipelineName string,
) (pipelineRun models.PipelineRunView, err error) {
  return this.
  compositionRoot.
  RunPipelineUseCase().
  Execute(
    pipelineName,
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
