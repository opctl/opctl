package models

func NewSetOpDescriptionReq(
pathToOp string,
description string,
) *SetOpDescriptionReq {

  return &SetOpDescriptionReq{
    PathToOp:pathToOp,
    Description:description,
  }

}

type SetOpDescriptionReq struct {
  PathToOp string
  Description  string
}
