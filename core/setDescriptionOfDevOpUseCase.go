package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type setDescriptionOfDevOpUseCase interface {
  Execute(
  req models.SetDescriptionOfDevOpReq,
  ) (err error)
}

func newSetDescriptionOfDevOpUseCase(
fs ports.Filesys,
yml yamlCodec,
) setDescriptionOfDevOpUseCase {

  return &_setDescriptionOfDevOpUseCase{
    fs:fs,
    yml:yml,
  }

}

type _setDescriptionOfDevOpUseCase struct {
  fs  ports.Filesys
  yml yamlCodec
}

func (this _setDescriptionOfDevOpUseCase) Execute(
req models.SetDescriptionOfDevOpReq,
) (err error) {

  var devOpFileBytes []byte
  devOpFileBytes, err = this.fs.ReadDevOpFile(req.DevOpName)
  if (nil != err) {
    return
  }

  devOpFile := devOpFile{}
  err = this.yml.fromYaml(
    devOpFileBytes,
    &devOpFile,
  )
  if (nil != err) {
    return
  }

  devOpFile.Description = req.Description

  devOpFileBytes, err = this.yml.toYaml(&devOpFile)
  if (nil != err) {
    return
  }

  err = this.fs.SaveDevOpFile(
    req.DevOpName,
    devOpFileBytes,
  )

  return

}
