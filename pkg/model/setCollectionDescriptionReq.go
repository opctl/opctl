package models

func NewSetCollectionDescriptionReq(
pathToCollection string,
description string,
) *SetCollectionDescriptionReq {

  return &SetCollectionDescriptionReq{
    PathToCollection:pathToCollection,
    Description:description,
  }

}

type SetCollectionDescriptionReq struct {
  PathToCollection string
  Description  string
}
