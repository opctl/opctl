package models

func NewAddOpReq(
projectUrl *Url,
name string,
description string,
) *AddOpReq {

  return &AddOpReq{
    ProjectUrl:projectUrl,
    Name:name,
    Description :description,
  }

}

type AddOpReq struct {
  ProjectUrl  *Url
  Name        string `json:"name"`
  Description string `json:"description"`
}

