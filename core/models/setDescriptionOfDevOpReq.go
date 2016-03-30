package models

func NewSetDescriptionOfDevOpReq(
projectUrl *ProjectUrl,
description string,
devOpName string,
) *SetDescriptionOfDevOpReq {

  return &SetDescriptionOfDevOpReq{
    ProjectUrl:projectUrl,
    Description:description,
    DevOpName :devOpName,
  }

}

type SetDescriptionOfDevOpReq struct {
  ProjectUrl  *ProjectUrl
  Description string
  DevOpName   string
}
