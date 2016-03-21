package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type setDescriptionOfDevOpUcExecuter interface {
  Execute(
  req models.SetDescriptionOfDevOpReq,
  ) (err error)
}

func newSetDescriptionOfDevOpUcExecuter(
fs ports.Filesys,
yml yamlCodec,
) setDescriptionOfDevOpUcExecuter {

  return &setDescriptionOfDevOpUcExecuterImpl{
    fs:fs,
    yml:yml,
  }

}

type setDescriptionOfDevOpUcExecuterImpl struct {
  fs  ports.Filesys
  yml yamlCodec
}

func (x setDescriptionOfDevOpUcExecuterImpl) Execute(
req models.SetDescriptionOfDevOpReq,
) (err error) {

  var devOpFileBytes []byte
  devOpFileBytes, err = x.fs.ReadDevOpFile(req.DevOpName)
  if (nil != err) {
    return
  }

  devOpFile := devOpFile{}
  err = x.yml.fromYaml(
    devOpFileBytes,
    &devOpFile,
  )
  if (nil != err) {
    return
  }

  devOpFile.Description = req.Description

  devOpFileBytes, err = x.yml.toYaml(&devOpFile)
  if (nil != err) {
    return
  }

  err = x.fs.SaveDevOpFile(
    req.DevOpName,
    devOpFileBytes,
  )

  return

}
