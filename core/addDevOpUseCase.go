package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type addDevOpUseCase interface {
  Execute(
  req models.AddDevOpReq,
  ) (err error)
}

func newAddDevOpUseCase(
fs ports.Filesys,
yml yamlCodec,
ce ports.ContainerEngine,
) addDevOpUseCase {

  return &_addDevOpUseCase{
    fs:fs,
    yml:yml,
    ce:ce,
  }

}

type _addDevOpUseCase struct {
  fs  ports.Filesys
  yml yamlCodec
  ce  ports.ContainerEngine
}

func (this _addDevOpUseCase) Execute(
req models.AddDevOpReq,
) (err error) {

  err = this.fs.CreateDevOpDir(req.Name)
  if (nil != err) {
    return
  }

  var devOpFile = devOpFile{Description:req.Description}

  var devOpFileBytes []byte
  devOpFileBytes, err = this.yml.toYaml(&devOpFile)
  if (nil != err) {
    return
  }

  err = this.fs.SaveDevOpFile(
    req.Name,
    devOpFileBytes,
  )
  if (nil != err) {
    return
  }

  err = this.ce.InitDevOp(req.Name)

  return

}

