package models

func NewAddPipelineReq(
name string,
description string,
) *AddPipelineReq {

  return &AddPipelineReq{
    Name:name,
    Description :description,
  }

}

type AddPipelineReq struct {
  Name        string
  Description string
}

