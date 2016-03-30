package models

func NewAddPipelineReq(
projectUrl *ProjectUrl,
name string,
description string,
) *AddPipelineReq {

  return &AddPipelineReq{
    ProjectUrl:projectUrl,
    Name:name,
    Description :description,
  }

}

type AddPipelineReq struct {
  ProjectUrl  *ProjectUrl
  Name        string
  Description string
}

