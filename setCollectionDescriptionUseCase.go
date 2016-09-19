package opspec

//go:generate counterfeiter -o ./fakeSetCollectionDescriptionUseCase.go --fake-name fakeSetCollectionDescriptionUseCase ./ setCollectionDescriptionUseCase

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
)

type setCollectionDescriptionUseCase interface {
  Execute(
  req models.SetCollectionDescriptionReq,
  ) (err error)
}

func newSetCollectionDescriptionUseCase(
filesystem filesystem,
yaml format,
) setCollectionDescriptionUseCase {

  return &_setCollectionDescriptionUseCase{
    filesystem:filesystem,
    yaml:yaml,
  }

}

type _setCollectionDescriptionUseCase struct {
  filesystem filesystem
  yaml       format
}

func (this _setCollectionDescriptionUseCase) Execute(
req models.SetCollectionDescriptionReq,
) (err error) {

  pathToCollectionManifest := path.Join(req.PathToCollection, NameOfCollectionManifestFile)

  collectionManifestBytes, err := this.filesystem.GetBytesOfFile(
    pathToCollectionManifest,
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

  collectionManifest.Description = req.Description

  collectionManifestBytes, err = this.yaml.From(&collectionManifest)
  if (nil != err) {
    return
  }

  err = this.filesystem.SaveFile(
    pathToCollectionManifest,
    collectionManifestBytes,
  )

  return

}
