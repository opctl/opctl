package models

func NewSetDescriptionOfDevOpReq(
pathToProjectRootDir string,
description string,
devOpName string,
) *SetDescriptionOfDevOpReq {

  return &SetDescriptionOfDevOpReq{
    PathToProjectRootDir:pathToProjectRootDir,
    Description:description,
    DevOpName :devOpName,
  }

}

type SetDescriptionOfDevOpReq struct {
  PathToProjectRootDir string
  Description          string
  DevOpName            string
}
