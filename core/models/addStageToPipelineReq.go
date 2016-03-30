package models

func NewAddStageToPipelineReq(
projectUrl *ProjectUrl,
isPipelineStage bool,
stageName string,
pipelineName string,
precedingStageName string,
) *AddStageToPipelineReq {

  return &AddStageToPipelineReq{
    ProjectUrl:projectUrl,
    IsPipelineStage:isPipelineStage,
    StageName :stageName,
    PipelineName :pipelineName,
    PrecedingStageName :precedingStageName,
  }

}

type AddStageToPipelineReq struct {
  ProjectUrl         *ProjectUrl
  IsPipelineStage    bool
  StageName          string
  PipelineName       string
  PrecedingStageName string
}