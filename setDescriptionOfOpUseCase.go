package sdk_golang

import "github.com/opspec-io/sdk-golang/models"

//go:generate counterfeiter -o ./fakeSetDescriptionOfOpUseCase.go --fake-name fakeSetDescriptionOfOpUseCase ./ setDescriptionOfOpUseCase

type setDescriptionOfOpUseCase interface {
  Execute(
  req models.SetDescriptionOfOpReq,
  ) (err error)
}

func newSetDescriptionOfOpUseCase(
filesystem Filesystem,
yamlCodec yamlCodec,
) setDescriptionOfOpUseCase {

  return &_setDescriptionOfOpUseCase{
    filesystem:filesystem,
    yamlCodec:yamlCodec,
  }

}

type _setDescriptionOfOpUseCase struct {
  filesystem Filesystem
  yamlCodec  yamlCodec
}

func (this _setDescriptionOfOpUseCase) Execute(
req models.SetDescriptionOfOpReq,
) (err error) {

  opFileBytes, err := this.filesystem.GetBytesOfFile(req.PathToOpFile)
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
    req.PathToOpFile,
    opFileBytes,
  )

  return

}
