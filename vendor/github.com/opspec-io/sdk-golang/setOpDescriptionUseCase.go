package opspec

//go:generate counterfeiter -o ./fakeSetOpDescriptionUseCase.go --fake-name fakeSetOpDescriptionUseCase ./ setOpDescriptionUseCase

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
)

type setOpDescriptionUseCase interface {
  Execute(
  req models.SetOpDescriptionReq,
  ) (err error)
}

func newSetOpDescriptionUseCase(
filesystem Filesystem,
yamlCodec yamlCodec,
) setOpDescriptionUseCase {

  return &_setOpDescriptionUseCase{
    filesystem:filesystem,
    yamlCodec:yamlCodec,
  }

}

type _setOpDescriptionUseCase struct {
  filesystem Filesystem
  yamlCodec  yamlCodec
}

func (this _setOpDescriptionUseCase) Execute(
req models.SetOpDescriptionReq,
) (err error) {

  pathToOpManifest := path.Join(req.PathToOp, NameOfOpManifestFile)

  opBytes, err := this.filesystem.GetBytesOfFile(
    pathToOpManifest,
  )
  if (nil != err) {
    return
  }

  opManifest := models.OpManifest{}
  err = this.yamlCodec.FromYaml(
    opBytes,
    &opManifest,
  )
  if (nil != err) {
    return
  }

  opManifest.Description = req.Description

  opBytes, err = this.yamlCodec.ToYaml(&opManifest)
  if (nil != err) {
    return
  }

  err = this.filesystem.SaveFile(
    pathToOpManifest,
    opBytes,
  )

  return

}
