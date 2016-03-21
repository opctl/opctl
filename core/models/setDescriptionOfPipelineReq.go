package models

func NewSetDescriptionOfPipelineReq(
description string,
pipelineName string,
) *SetDescriptionOfPipelineReq {

  return &SetDescriptionOfPipelineReq{
    Description:description,
    PipelineName :pipelineName,
  }

}

type SetDescriptionOfPipelineReq struct {
  Description  string
  PipelineName string
}
