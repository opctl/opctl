package models

func NewAddDevOpReq(
pathToProjectRootDir string,
name string,
description string,
) *AddDevOpReq {

  return &AddDevOpReq{
    PathToProjectRootDir:pathToProjectRootDir,
    Name:name,
    Description :description,
  }

}

type AddDevOpReq struct {
  PathToProjectRootDir string
  Name                 string
  Description          string
}

