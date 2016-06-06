package models

func NewAddOpReq(
path string,
name string,
description string,
) *AddOpReq {

  return &AddOpReq{
    Path:path,
    Name:name,
    Description :description,
  }

}

type AddOpReq struct {
  Path        string
  Name        string
  Description string
}

