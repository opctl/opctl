package models

func NewAddPipelineReq(
pathToProjectRootDir string,
name string,
description string,
) *AddPipelineReq {

  return &AddPipelineReq{
    PathToProjectRootDir:pathToProjectRootDir,
    Name:name,
    Description :description,
  }

}

type AddPipelineReq struct {
  PathToProjectRootDir string
  Name                 string
  Description          string
}

