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

  pathToOpFile := path.Join(req.PathToOp, NameOfOpFile)

  opBytes, err := this.filesystem.GetBytesOfFile(
    pathToOpFile,
  )
  if (nil != err) {
    return
  }

  opFile := models.OpFile{}
  err = this.yamlCodec.FromYaml(
    opBytes,
    &opFile,
  )
  if (nil != err) {
    return
  }

  opFile.Description = req.Description

  opBytes, err = this.yamlCodec.ToYaml(&opFile)
  if (nil != err) {
    return
  }

  err = this.filesystem.SaveFile(
    pathToOpFile,
    opBytes,
  )

  return

}
