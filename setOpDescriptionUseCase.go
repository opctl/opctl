package opspec

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
)

//go:generate counterfeiter -o ./fakeSetOpDescriptionUseCase.go --fake-name fakeSetOpDescriptionUseCase ./ setOpDescriptionUseCase

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

  opBytes, err := this.filesystem.GetBytesOfFile(
    path.Join(req.PathToOp, NameOfOpFile),
  )
  if (nil != err) {
    return
  }

  opFile := models.OpFile{}
  err = this.yamlCodec.fromYaml(
    opBytes,
    &opFile,
  )
  if (nil != err) {
    return
  }

  opFile.Description = req.Description

  opBytes, err = this.yamlCodec.toYaml(&opFile)
  if (nil != err) {
    return
  }

  err = this.filesystem.SaveFile(
    req.PathToOp,
    opBytes,
  )

  return

}
