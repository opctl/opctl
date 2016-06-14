package models

type OpSummaryView struct {
  Name string
}

func NewOpSummaryView(
name string,
) *OpSummaryView {

  return &OpSummaryView{
    Name:name,
  }

}
