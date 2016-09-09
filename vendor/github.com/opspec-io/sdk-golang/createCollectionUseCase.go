package opspec

//go:generate counterfeiter -o ./fakeCreateCollectionUseCase.go --fake-name fakeCreateCollectionUseCase ./ createCollectionUseCase

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
)

type createCollectionUseCase interface {
  Execute(
  req models.CreateCollectionReq,
  ) (err error)
}

func newCreateCollectionUseCase(
filesystem Filesystem,
yamlCodec yamlCodec,
) createCollectionUseCase {

  return &_createCollectionUseCase{
    filesystem:filesystem,
    yamlCodec:yamlCodec,
  }

}

type _createCollectionUseCase struct {
  filesystem Filesystem
  yamlCodec  yamlCodec
}

func (this _createCollectionUseCase) Execute(
req models.CreateCollectionReq,
) (err error) {

  err = this.filesystem.AddDir(
    req.Path,
  )
  if (nil != err) {
    return
  }

  var opCollection = models.CollectionManifest{
    Manifest:models.Manifest{
      Description:req.Description,
      Name:req.Name,
    },
  }

  opManifestBytes, err := this.yamlCodec.ToYaml(&opCollection)
  if (nil != err) {
    return
  }

  err = this.filesystem.SaveFile(
    path.Join(req.Path, "collection.yml"),
    opManifestBytes,
  )

  return

}
