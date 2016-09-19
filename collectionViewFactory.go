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
filesystem filesystem,
opViewFactory opViewFactory,
yaml format,
) collectionViewFactory {

  return &_collectionViewFactory{
    filesystem:filesystem,
    opViewFactory:opViewFactory,
    yaml:yaml,
  }

}

type _collectionViewFactory struct {
  filesystem    filesystem
  opViewFactory opViewFactory
  yaml          format
}

func (this _collectionViewFactory) Construct(
collectionBundlePath string,
) (
collectionView models.CollectionView,
err error,
) {

  collectionManifestPath := path.Join(collectionBundlePath, NameOfCollectionManifestFile)

  collectionManifestBytes, err := this.filesystem.GetBytesOfFile(
    collectionManifestPath,
  )
  if (nil != err) {
    return
  }

  collectionManifest := models.CollectionManifest{}
  err = this.yaml.To(
    collectionManifestBytes,
    &collectionManifest,
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
    collectionManifest.Description,
    collectionManifest.Name,
    opViews,
  )

  return

}

