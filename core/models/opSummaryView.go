package models

type OpSummaryView struct {
  Name string `json:"name"`
}

func NewOpSummaryView(
name string,
) *OpSummaryView {

  return &OpSummaryView{
    Name:name,
  }

}