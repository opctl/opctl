package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type addPipelineUseCase interface {
  Execute(
  req models.AddPipelineReq,
  ) (err error)
}

func newAddPipelineUseCase(
fs ports.Filesys,
yml yamlCodec,
) addPipelineUseCase {

  return &_addPipelineUseCase{
    fs:fs,
    yml:yml,
  }

}

type _addPipelineUseCase struct {
  fs  ports.Filesys
  yml yamlCodec
}

func (this _addPipelineUseCase) Execute(
req models.AddPipelineReq,
) (err error) {

  err = this.fs.CreatePipelineDir(req.Name)
  if (nil != err) {
    return
  }

  var pipelineFile = pipelineFile{Description:req.Description}

  var pipelineFileBytes []byte
  pipelineFileBytes, err = this.yml.toYaml(&pipelineFile)
  if (nil != err) {
    return
  }

  err = this.fs.SavePipelineFile(
    req.Name,
    pipelineFileBytes,
  )

  return

}
