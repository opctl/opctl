package models

func NewSetDescriptionOfOpReq(
pathToOpFile string,
description string,
) *SetDescriptionOfOpReq {

  return &SetDescriptionOfOpReq{
    PathToOpFile:pathToOpFile,
    Description:description,
  }

}

type SetDescriptionOfOpReq struct {
  PathToOpFile string
  Description  string
}
