package models

func NewOpView(
description string,
inputs []Param,
name string,
run *RunDeclaration,
version string,
) *OpView {

  return &OpView{
    Description:description,
    Inputs:inputs,
    Name:name,
    Run:run,
    Version:version,
  }

}

type OpView struct {
  Description string
  Inputs      []Param
  Name        string
  Run         *RunDeclaration
  Version     string
}
