package opspec

//go:generate counterfeiter -o ./fakeTryResolveDefaultCollectionUseCase.go --fake-name fakeTryResolveDefaultCollectionUseCase ./ tryResolveDefaultCollectionUseCase

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
  "path/filepath"
  "strings"
)

type tryResolveDefaultCollectionUseCase interface {
  Execute(
  req models.TryResolveDefaultCollectionReq,
  ) (
  pathToDefaultCollection string,
  err error,
  )
}

func newTryResolveDefaultCollectionUseCase(
filesystem filesystem,
) tryResolveDefaultCollectionUseCase {

  return &_tryResolveDefaultCollectionUseCase{
    filesystem:filesystem,
  }

}

type _tryResolveDefaultCollectionUseCase struct {
  filesystem filesystem
}

func (this _tryResolveDefaultCollectionUseCase) Execute(
req models.TryResolveDefaultCollectionReq,
) (
pathToDefaultCollection string,
err error,
) {

  pathToCurrentDir, err := filepath.Abs(req.PathToDir)
  if (nil != err) {
    return
  }

  childDirFileInfos, err := this.filesystem.ListChildFileInfosOfDir(pathToCurrentDir)
  if (nil != err) {
    return
  }

  for _, x := range childDirFileInfos {

    if x.IsDir() && x.Name() == NameOfDefaultOpCollectionDir {
      // handle default collection found
      pathToDefaultCollection = path.Join(pathToCurrentDir, x.Name())
      return
    }

  }

  if (!strings.HasSuffix(filepath.Clean(pathToCurrentDir), "/")) {
    // handle non root path
    pathToParentDir := filepath.Dir(pathToCurrentDir)
    return this.Execute(
      *models.NewTryResolveDefaultCollectionReq(pathToParentDir),
    )
  }

  return

}
