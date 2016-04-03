package models

func NewSetDescriptionOfOperationReq(
projectUrl *ProjectUrl,
description string,
operationName string,
) *SetDescriptionOfOperationReq {

  return &SetDescriptionOfOperationReq{
    ProjectUrl:projectUrl,
    Description:description,
    OperationName :operationName,
  }

}

type SetDescriptionOfOperationReq struct {
  ProjectUrl    *ProjectUrl
  Description   string
  OperationName string
}
