package models

func NewSetDescriptionOfPipelineReq(
projectUrl *ProjectUrl,
description string,
pipelineName string,
) *SetDescriptionOfPipelineReq {

  return &SetDescriptionOfPipelineReq{
    ProjectUrl:projectUrl,
    Description:description,
    PipelineName :pipelineName,
  }

}

type SetDescriptionOfPipelineReq struct {
  ProjectUrl   *ProjectUrl
  Description  string
  PipelineName string
}
