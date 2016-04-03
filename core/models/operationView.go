package models

type OperationView struct {
  Description   string
  Name          string
  SubOperations []SubOperationView `json:",omitempty"`
}

func NewOperationView(
description string,
name string,
subOperations []SubOperationView,
) *OperationView {

  return &OperationView{
    Description:description,
    Name:name,
    SubOperations:subOperations,
  }

}