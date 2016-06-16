package models

func NewOpParamView(
name string,
description string,
isSecret bool,
) *OpParamView {

  return &OpParamView{
    Name:name,
    Description:description,
    IsSecret:isSecret,
  }

}

type OpParamView struct {
  Name        string
  Description string
  IsSecret    bool
}
