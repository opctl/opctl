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
filesystem Filesystem,
yamlCodec yamlCodec,
) setCollectionDescriptionUseCase {

  return &_setCollectionDescriptionUseCase{
    filesystem:filesystem,
    yamlCodec:yamlCodec,
  }

}

type _setCollectionDescriptionUseCase struct {
  filesystem Filesystem
  yamlCodec  yamlCodec
}

func (this _setCollectionDescriptionUseCase) Execute(
req models.SetCollectionDescriptionReq,
) (err error) {

  pathToCollectionBundleManifest := path.Join(req.PathToCollection, NameOfCollectionBundleManifest)

  collectionBundleManifestBytes, err := this.filesystem.GetBytesOfFile(
    pathToCollectionBundleManifest,
  )
  if (nil != err) {
    return
  }

  collectionBundleManifest := models.CollectionBundleManifest{}
  err = this.yamlCodec.FromYaml(
    collectionBundleManifestBytes,
    &collectionBundleManifest,
  )
  if (nil != err) {
    return
  }

  collectionBundleManifest.Description = req.Description

  collectionBundleManifestBytes, err = this.yamlCodec.ToYaml(&collectionBundleManifest)
  if (nil != err) {
    return
  }

  err = this.filesystem.SaveFile(
    pathToCollectionBundleManifest,
    collectionBundleManifestBytes,
  )

  return

}
