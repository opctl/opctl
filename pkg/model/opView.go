package model

type OpView struct {
  Description string
  Inputs      []Param
  Name        string
  Run         *RunDeclaration
  Version     string
}
