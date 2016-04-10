package models

func NewSetDescriptionOfOpReq(
projectUrl *Url,
opName string,
description string,
) *SetDescriptionOfOpReq {

  return &SetDescriptionOfOpReq{
    ProjectUrl:projectUrl,
    OpName:opName,
    Description:description,
  }

}

type SetDescriptionOfOpReq struct {
  ProjectUrl  *Url
  OpName      string
  Description string `json:"description"`
}
