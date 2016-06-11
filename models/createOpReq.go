package models

func NewCreateOpReq(
path string,
name string,
description string,
) *CreateOpReq {

  return &CreateOpReq{
    Path:path,
    Name:name,
    Description :description,
  }

}

type CreateOpReq struct {
  Path        string
  Name        string
  Description string
}

