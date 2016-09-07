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

  pathToOpBundleManifest := path.Join(req.PathToOp, NameOfOpBundleManifest)

  opBytes, err := this.filesystem.GetBytesOfFile(
    pathToOpBundleManifest,
  )
  if (nil != err) {
    return
  }

  opBundleManifest := models.OpBundleManifest{}
  err = this.yamlCodec.FromYaml(
    opBytes,
    &opBundleManifest,
  )
  if (nil != err) {
    return
  }

  opBundleManifest.Description = req.Description

  opBytes, err = this.yamlCodec.ToYaml(&opBundleManifest)
  if (nil != err) {
    return
  }

  err = this.filesystem.SaveFile(
    pathToOpBundleManifest,
    opBytes,
  )

  return

}
