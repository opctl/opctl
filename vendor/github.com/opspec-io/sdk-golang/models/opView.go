package models

func NewOpView(
description string,
inputs []Parameter,
name string,
outputs []Parameter,
run RunInstruction,
version string,
) *OpView {

  return &OpView{
    Description:description,
    Inputs:inputs,
    Name:name,
    Outputs:outputs,
    Run:run,
    Version:version,
  }

}

type OpView struct {
  Description string
  Inputs      []Parameter
  Name        string
  Outputs     []Parameter
  Run         RunInstruction
  Version     string
}
