package core

import (
  "github.com/dev-op-spec/engine/core/models"
"github.com/dev-op-spec/engine/core/ports"
)

type addDevOpUcExecuter interface {
  Execute(
  req models.AddDevOpReq,
  ) (err error)
}

func newAddDevOpUcExecuter(
fs ports.Filesys,
yml yamlCodec,
ce ports.ContainerEngine,
) addDevOpUcExecuter {

  return &addDevOpUcExecuterImpl{
    fs:fs,
    yml:yml,
    ce:ce,
  }

}

type addDevOpUcExecuterImpl struct {
  fs  ports.Filesys
  yml yamlCodec
  ce  ports.ContainerEngine
}

func (x addDevOpUcExecuterImpl) Execute(
req models.AddDevOpReq,
) (err error) {

  err = x.fs.CreateDevOpDir(req.Name)
  if (nil != err) {
    return
  }

  var devOpFile = devOpFile{Description:req.Description}

  var devOpFileBytes []byte
  devOpFileBytes, err = x.yml.toYaml(&devOpFile)
  if (nil != err) {
    return
  }

  err = x.fs.SaveDevOpFile(
    req.Name,
    devOpFileBytes,
  )
  if (nil != err) {
    return
  }

  err = x.ce.InitDevOp(req.Name)

  return

}

