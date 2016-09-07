package models

func NewOpView(
description string,
inputs []Param,
name string,
run *RunStatement,
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
  Run         *RunStatement
  Version     string
}
