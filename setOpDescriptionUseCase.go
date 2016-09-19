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
filesystem filesystem,
yaml format,
) setOpDescriptionUseCase {

  return &_setOpDescriptionUseCase{
    filesystem:filesystem,
    yaml:yaml,
  }

}

type _setOpDescriptionUseCase struct {
  filesystem filesystem
  yaml       format
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
  err = this.yaml.To(
    opBytes,
    &opManifest,
  )
  if (nil != err) {
    return
  }

  opManifest.Description = req.Description

  opBytes, err = this.yaml.From(&opManifest)
  if (nil != err) {
    return
  }

  err = this.filesystem.SaveFile(
    pathToOpManifest,
    opBytes,
  )

  return

}
