package models

func NewAddDevOpReq(
projectUrl *ProjectUrl,
name string,
description string,
) *AddDevOpReq {

  return &AddDevOpReq{
    ProjectUrl:projectUrl,
    Name:name,
    Description :description,
  }

}

type AddDevOpReq struct {
  ProjectUrl  *ProjectUrl
  Name        string
  Description string
}

