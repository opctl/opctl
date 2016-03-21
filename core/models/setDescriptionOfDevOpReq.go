package models

func NewSetDescriptionOfDevOpReq(
description string,
devOpName string,
) *SetDescriptionOfDevOpReq {

  return &SetDescriptionOfDevOpReq{
    Description:description,
    DevOpName :devOpName,
  }

}

type SetDescriptionOfDevOpReq struct {
  Description string
  DevOpName   string
}
