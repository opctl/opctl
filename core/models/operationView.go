package models

type OperationView struct {
  Description   string
  Name          string
  SubOperations []OperationRefView `json:",omitempty"`
}

func NewOperationView(
description string,
name string,
subOperations []OperationRefView,
) *OperationView {

  return &OperationView{
    Description:description,
    Name:name,
    SubOperations:subOperations,
  }

}