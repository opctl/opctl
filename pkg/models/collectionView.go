package models

func NewCollectionView(
description string,
name string,
ops []OpView,
) *CollectionView {

  return &CollectionView{
    Description:description,
    Name:name,
    Ops:ops,
  }

}

type CollectionView struct {
  Description string
  Name        string
  Ops         []OpView
}
