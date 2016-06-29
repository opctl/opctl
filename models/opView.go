package models

func NewOpView(
description string,
name string,
params []OpParamView,
subOps []SubOpView,
) *OpView {

  return &OpView{
    Description:description,
    Name:name,
    Params:params,
    SubOps:subOps,
  }

}

type OpView struct {
  Description string
  Name        string
  Params      []OpParamView
  SubOps      []SubOpView
}
