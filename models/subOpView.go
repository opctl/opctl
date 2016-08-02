package models

func NewSubOpView(
isParallel bool,
url string,
) *SubOpView {

  return &SubOpView{
    IsParallel:isParallel,
    Url:url,
  }

}

type SubOpView struct {
  IsParallel bool
  Url        string
}


