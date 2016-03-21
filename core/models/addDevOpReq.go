package models

func NewAddDevOpReq(
name string,
description string,
) *AddDevOpReq {

  return &AddDevOpReq{
    Name:name,
    Description :description,
  }

}

type AddDevOpReq struct {
  Name        string
  Description string
}

