package models

type SubOpView struct {
  IsParallel bool
  Url        string
}

func NewSubOpView(
isParallel bool,
url string,
) *SubOpView {

  return &SubOpView{
    IsParallel:isParallel,
    Url:url,
  }

}
