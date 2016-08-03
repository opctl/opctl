package models

func NewCreateCollectionReq(
path string,
name string,
description string,
) *CreateCollectionReq {

  return &CreateCollectionReq{
    Path:path,
    Name:name,
    Description :description,
  }

}

type CreateCollectionReq struct {
  Path        string
  Name        string
  Description string
}

