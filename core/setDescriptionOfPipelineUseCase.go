package core

import (
  "github.com/dev-op-spec/engine/core/ports"
  "github.com/dev-op-spec/engine/core/models"
)

type setDescriptionOfPipelineUseCase interface {
  Execute(
  req models.SetDescriptionOfPipelineReq,
  ) (err error)
}

func newSetDescriptionOfPipelineUseCase(
fs ports.Filesys,
yml yamlCodec,
) setDescriptionOfPipelineUseCase {

  return &_setDescriptionOfPipelineUseCase{
    fs:fs,
    yml:yml,
  }

}

type _setDescriptionOfPipelineUseCase struct {
  fs  ports.Filesys
  yml yamlCodec
}

func (this _setDescriptionOfPipelineUseCase) Execute(
req models.SetDescriptionOfPipelineReq,
) (err error) {

  var pipelineFileBytes []byte
  pipelineFileBytes, err = this.fs.ReadPipelineFile(req.PipelineName)
  if (nil != err) {
    return
  }

  pipelineFile := pipelineFile{}
  err = this.yml.fromYaml(
    pipelineFileBytes,
    &pipelineFile,
  )
  if (nil != err) {
    return
  }

  pipelineFile.Description = req.Description

  pipelineFileBytes, err = this.yml.toYaml(&pipelineFile)
  if (nil != err) {
    return
  }

  err = this.fs.SavePipelineFile(
    req.PipelineName,
    pipelineFileBytes,
  )

  return

}
