package models

func NewAddStageToPipelineReq(
isPipelineStage bool,
stageName string,
pipelineName string,
precedingStageName string,
) *AddStageToPipelineReq {

  return &AddStageToPipelineReq{
    IsPipelineStage:isPipelineStage,
    StageName :stageName,
    PipelineName :pipelineName,
    PrecedingStageName :precedingStageName,
  }

}

type AddStageToPipelineReq struct {
  IsPipelineStage    bool
  StageName          string
  PipelineName       string
  PrecedingStageName string
}