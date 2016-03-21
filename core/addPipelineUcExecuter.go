package core

import (
  "github.com/dev-op-spec/engine/core/models"
"github.com/dev-op-spec/engine/core/ports"
)

type addPipelineUcExecuter interface {
  Execute(
  req models.AddPipelineReq,
  ) (err error)
}

func newAddPipelineUcExecuter(
fs ports.Filesys,
yml yamlCodec,
) addPipelineUcExecuter {

  return &addPipelineUcExecuterImpl{
    fs:fs,
    yml:yml,
  }

}

type addPipelineUcExecuterImpl struct {
  fs  ports.Filesys
  yml yamlCodec
}

func (x addPipelineUcExecuterImpl) Execute(
req models.AddPipelineReq,
) (err error) {

  err = x.fs.CreatePipelineDir(req.Name)
  if (nil != err) {
    return
  }

  var pipelineFile = pipelineFile{Description:req.Description}

  var pipelineFileBytes []byte
  pipelineFileBytes, err = x.yml.toYaml(&pipelineFile)
  if (nil != err) {
    return
  }

  err = x.fs.SavePipelineFile(
    req.Name,
    pipelineFileBytes,
  )

  return

}
