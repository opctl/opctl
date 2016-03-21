package core

import (
  "github.com/dev-op-spec/engine/core/ports"
  "github.com/dev-op-spec/engine/core/models"
)

type setDescriptionOfPipelineUcExecuter interface {
  Execute(
  req models.SetDescriptionOfPipelineReq,
  ) (err error)
}

func newSetDescriptionOfPipelineUcExecuter(
fs ports.Filesys,
yml yamlCodec,
) setDescriptionOfPipelineUcExecuter {

  return &setDescriptionOfPipelineUcExecuterImpl{
    fs:fs,
    yml:yml,
  }

}

type setDescriptionOfPipelineUcExecuterImpl struct {
  fs  ports.Filesys
  yml yamlCodec
}

func (x setDescriptionOfPipelineUcExecuterImpl) Execute(
req models.SetDescriptionOfPipelineReq,
) (err error) {

  var pipelineFileBytes []byte
  pipelineFileBytes, err = x.fs.ReadPipelineFile(req.PipelineName)
  if (nil != err) {
    return
  }

  pipelineFile := pipelineFile{}
  err = x.yml.fromYaml(
    pipelineFileBytes,
    &pipelineFile,
  )
  if (nil != err) {
    return
  }

  pipelineFile.Description = req.Description

  pipelineFileBytes, err = x.yml.toYaml(&pipelineFile)
  if (nil != err) {
    return
  }

  err = x.fs.SavePipelineFile(
    req.PipelineName,
    pipelineFileBytes,
  )

  return

}
