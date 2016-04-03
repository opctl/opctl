package models

type OperationSummaryView struct {
  Name string `json:"name"`
}

func NewOperationSummaryView(
name string,
) *OperationSummaryView {

  return &OperationSummaryView{
    Name:name,
  }

}