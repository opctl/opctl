package model

func NewTryResolveDefaultCollectionReq(
pathToDir string,
) *TryResolveDefaultCollectionReq {

  return &TryResolveDefaultCollectionReq{
    PathToDir:pathToDir,
  }

}

type TryResolveDefaultCollectionReq struct {
  PathToDir string
}
