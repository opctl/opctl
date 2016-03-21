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

func (_api _api) AddDevOp(
req models.AddDevOpReq,
) (err error) {
  return _api.
  compositionRoot.
  AddDevOpUcExecuter().
  Execute(req)
}

func (_api _api) AddPipeline(
req models.AddPipelineReq,
) (err error) {
  return _api.
  compositionRoot.
  AddPipelineUcExecuter().
  Execute(req)
}

func (_api _api) AddStageToPipeline(
req models.AddStageToPipelineReq,
) (err error) {
  return _api.
  compositionRoot.
  AddStageToPipelineUcExecuter().
  Execute(req)
}

func (_api _api) ListDevOps(
) (devOps []models.DevOpView, err error) {
  return _api.
  compositionRoot.
  ListDevOpsUcExecuter().
  Execute()
}

func (_api _api) ListPipelines(
) (pipelines []models.PipelineView, err error) {
  return _api.
  compositionRoot.
  ListPipelinesUcExecuter().
  Execute()
}

func (_api _api) RunDevOp(
devOpName string,
) (devOpRun models.DevOpRunView, err error) {
  return _api.
  compositionRoot.
  RunDevOpUcExecuter().
  Execute(devOpName)
}

func (_api _api) RunPipeline(
pipelineName string,
) (pipelineRun models.PipelineRunView, err error) {
  return _api.
  compositionRoot.
  RunPipelineUcExecuter().
  Execute(
    pipelineName,
    make([]string, 0),
  )
}

func (_api _api) SetDescriptionOfDevOp(
req models.SetDescriptionOfDevOpReq,
) (err error) {
  return _api.
  compositionRoot.
  SetDescriptionOfDevOpUcExecuter().
  Execute(
    req,
  )
}

func (_api _api) SetDescriptionOfPipeline(
req models.SetDescriptionOfPipelineReq,
) (err error) {
  return _api.
  compositionRoot.
  SetDescriptionOfPipelineUcExecuter().
  Execute(
    req,
  )
}
