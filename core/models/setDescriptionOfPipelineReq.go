package models

func NewSetDescriptionOfPipelineReq(
pathToProjectRootDir string,
description string,
pipelineName string,
) *SetDescriptionOfPipelineReq {

  return &SetDescriptionOfPipelineReq{
    PathToProjectRootDir:pathToProjectRootDir,
    Description:description,
    PipelineName :pipelineName,
  }

}

type SetDescriptionOfPipelineReq struct {
  PathToProjectRootDir string
  Description          string
  PipelineName         string
}
