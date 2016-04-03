package models

func NewAddOperationReq(
projectUrl *ProjectUrl,
name string,
description string,
) *AddOperationReq {

  return &AddOperationReq{
    ProjectUrl:projectUrl,
    Name:name,
    Description :description,
  }

}

type AddOperationReq struct {
  ProjectUrl  *ProjectUrl
  Name        string `json:"name"`
  Description string `json:"description"`
}

