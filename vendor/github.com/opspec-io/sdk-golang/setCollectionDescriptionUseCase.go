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

  pathToCollectionFile := path.Join(req.PathToCollection, NameOfCollectionFile)

  collectionFileBytes, err := this.filesystem.GetBytesOfFile(
    pathToCollectionFile,
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

  collectionFile.Description = req.Description

  collectionFileBytes, err = this.yamlCodec.ToYaml(&collectionFile)
  if (nil != err) {
    return
  }

  err = this.filesystem.SaveFile(
    pathToCollectionFile,
    collectionFileBytes,
  )

  return

}
