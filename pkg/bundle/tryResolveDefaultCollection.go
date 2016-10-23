package bundle

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
  "path/filepath"
  "strings"
)

func (this _bundle) TryResolveDefaultCollection(
req models.TryResolveDefaultCollectionReq,
) (
pathToDefaultCollection string,
err error,
) {

  pathToCurrentDir, err := filepath.Abs(req.PathToDir)
  if (nil != err) {
    return
  }

  childDirFileInfos, err := this.fileSystem.ListChildFileInfosOfDir(pathToCurrentDir)
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
    return this.TryResolveDefaultCollection(
      *models.NewTryResolveDefaultCollectionReq(pathToParentDir),
    )
  }

  return

}
