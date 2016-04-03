package models

type OperationRefView struct {
  Name string
}

func NewOperationRefView(
name string,
) *OperationRefView {

  return &OperationRefView{
    Name:name,
  }

}