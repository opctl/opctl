package opspec

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
)

//go:generate counterfeiter -o ./fakeSetCollectionDescriptionUseCase.go --fake-name fakeSetCollectionDescriptionUseCase ./ setCollectionDescriptionUseCase

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

  opFileBytes, err := this.filesystem.GetBytesOfFile(
    path.Join(req.PathToCollection, NameOfCollectionFile),
  )
  if (nil != err) {
    return
  }

  opFile := models.OpFile{}
  err = this.yamlCodec.fromYaml(
    opFileBytes,
    &opFile,
  )
  if (nil != err) {
    return
  }

  opFile.Description = req.Description

  opFileBytes, err = this.yamlCodec.toYaml(&opFile)
  if (nil != err) {
    return
  }

  err = this.filesystem.SaveFile(
    req.PathToCollection,
    opFileBytes,
  )

  return

}
