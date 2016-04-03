package models

func NewOperationDetailedView(
description string,
name string,
subOperations []OperationSummaryView,
) *OperationDetailedView {

  return &OperationDetailedView{
    Description:description,
    Name:name,
    SubOperations:subOperations,
  }

}

type OperationDetailedView struct {
  Description   string `json:"description"`
  Name          string `json:"name"`
  SubOperations []OperationSummaryView `json:"subOperations,omitempty"`
}