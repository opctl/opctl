package opspec

//go:generate counterfeiter -o ./fakeCollectionViewFactory.go --fake-name fakeCollectionViewFactory ./ collectionViewFactory

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
)

type collectionViewFactory interface {
  Construct(
  collectionBundlePath string,
  ) (
  collectionView models.CollectionView,
  err error,
  )
}

func newCollectionViewFactory(
filesystem Filesystem,
opViewFactory opViewFactory,
yamlCodec yamlCodec,
) collectionViewFactory {

  return &_collectionViewFactory{
    filesystem:filesystem,
    opViewFactory:opViewFactory,
    yamlCodec:yamlCodec,
  }

}

type _collectionViewFactory struct {
  filesystem    Filesystem
  opViewFactory opViewFactory
  yamlCodec     yamlCodec
}

func (this _collectionViewFactory) Construct(
collectionBundlePath string,
) (
collectionView models.CollectionView,
err error,
) {

  collectionFilePath := path.Join(collectionBundlePath, NameOfCollectionFile)

  collectionFileBytes, err := this.filesystem.GetBytesOfFile(
    collectionFilePath,
  )
  if (nil != err) {
    return
  }

  collectionFile := models.CollectionFile{}
  err = this.yamlCodec.FromYaml(
    collectionFileBytes,
    &collectionFile,
  )
  if (nil != err) {
    return
  }

  var opViews []models.OpView
  childFileInfos, err := this.filesystem.ListChildFileInfosOfDir(collectionBundlePath)
  if (nil != err) {
    return
  }

  for _, childFileInfo := range childFileInfos {
    opView, err := this.opViewFactory.Construct(
      path.Join(collectionBundlePath, childFileInfo.Name()),
    )
    if (nil == err) {
      opViews = append(opViews, opView)
    }
  }

  collectionView = *models.NewCollectionView(
    collectionFile.Description,
    collectionFile.Name,
    opViews,
  )

  return

}

