package models

func NewAddOperationReq(
projectUrl *Url,
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
  ProjectUrl  *Url
  Name        string `json:"name"`
  Description string `json:"description"`
}

