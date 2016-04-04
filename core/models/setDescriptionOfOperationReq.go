package models

func NewSetDescriptionOfOperationReq(
projectUrl *Url,
operationName *string,
description string,
) *SetDescriptionOfOperationReq {

  return &SetDescriptionOfOperationReq{
    ProjectUrl:projectUrl,
    OperationName:operationName,
    Description:description,
  }

}

type SetDescriptionOfOperationReq struct {
  ProjectUrl    *Url
  OperationName *string
  Description   string `json:"description"`
}
