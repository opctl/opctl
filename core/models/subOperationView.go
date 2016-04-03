package models

type SubOperationView struct {
  Name string
}

func NewSubOperationView(
name string,
) *SubOperationView {

  return &SubOperationView{
    Name:name,
  }

}