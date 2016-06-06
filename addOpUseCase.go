package sdk

//go:generate counterfeiter -o ./fakeAddOpUseCase.go --fake-name fakeAddOpUseCase ./ addOpUseCase

import (
  "github.com/opspec-io/sdk-golang/models"
  "path"
)

type addOpUseCase interface {
  Execute(
  req models.AddOpReq,
  ) (err error)
}

func newAddOpUseCase(
filesystem Filesystem,
yamlCodec yamlCodec,
) addOpUseCase {

  return &_addOpUseCase{
    filesystem:filesystem,
    yamlCodec:yamlCodec,
  }

}

type _addOpUseCase struct {
  filesystem Filesystem
  yamlCodec  yamlCodec
}

func (this _addOpUseCase) Execute(
req models.AddOpReq,
) (err error) {

  err = this.filesystem.AddDir(
    req.Path,
  )
  if (nil != err) {
    return
  }

  var opFile = models.OpFile{
    Description:req.Description,
    Name:req.Name,
  }

  opFileBytes, err := this.yamlCodec.toYaml(&opFile)
  if (nil != err) {
    return
  }

  err = this.filesystem.SaveFile(
    path.Join(req.Path, "op.yml"),
    opFileBytes,
  )

  return

}
