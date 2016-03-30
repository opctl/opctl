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
filesys ports.Filesys,
pathToPipelineFileFactory pathToPipelineFileFactory,
yamlCodec yamlCodec,
) setDescriptionOfPipelineUseCase {

  return &_setDescriptionOfPipelineUseCase{
    filesys:filesys,
    pathToPipelineFileFactory:pathToPipelineFileFactory,
    yamlCodec:yamlCodec,
  }

}

type _setDescriptionOfPipelineUseCase struct {
  filesys                  ports.Filesys
  pathToPipelineFileFactory pathToPipelineFileFactory
  yamlCodec                yamlCodec
}

func (this _setDescriptionOfPipelineUseCase) Execute(
req models.SetDescriptionOfPipelineReq,
) (err error) {

  pathToPipelineFile := this.pathToPipelineFileFactory.Construct(
    req.ProjectUrl,
    req.PipelineName,
  )

  pipelineFileBytes, err := this.filesys.GetBytesOfFile(pathToPipelineFile)
  if (nil != err) {
    return
  }

  pipelineFile := pipelineFile{}
  err = this.yamlCodec.fromYaml(
    pipelineFileBytes,
    &pipelineFile,
  )
  if (nil != err) {
    return
  }

  pipelineFile.Description = req.Description

  pipelineFileBytes, err = this.yamlCodec.toYaml(&pipelineFile)
  if (nil != err) {
    return
  }

  err = this.filesys.SaveFile(
    pathToPipelineFile,
    pipelineFileBytes,
  )

  return

}
