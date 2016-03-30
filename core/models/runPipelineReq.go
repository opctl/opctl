package models

func NewRunPipelineReq(
projectUrl *ProjectUrl,
pipelineName string,
) *RunPipelineReq {

  return &RunPipelineReq{
    ProjectUrl:projectUrl,
    PipelineName :pipelineName,
  }

}

type RunPipelineReq struct {
  ProjectUrl   *ProjectUrl
  PipelineName string
}
