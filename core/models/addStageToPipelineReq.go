package models

func NewAddStageToPipelineReq(
pathToProjectRootDir string,
isPipelineStage bool,
stageName string,
pipelineName string,
precedingStageName string,
) *AddStageToPipelineReq {

  return &AddStageToPipelineReq{
    PathToProjectRootDir:pathToProjectRootDir,
    IsPipelineStage:isPipelineStage,
    StageName :stageName,
    PipelineName :pipelineName,
    PrecedingStageName :precedingStageName,
  }

}

type AddStageToPipelineReq struct {
  PathToProjectRootDir string
  IsPipelineStage      bool
  StageName            string
  PipelineName         string
  PrecedingStageName   string
}