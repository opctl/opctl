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
filesys ports.Filesys,
pathToPipelineDirFactory pathToPipelineFileFactory,
pathToPipelineFileFactory pathToPipelineFileFactory,
yamlCodec yamlCodec,
) addPipelineUseCase {

  return &_addPipelineUseCase{
    filesys:filesys,
    pathToPipelineDirFactory:pathToPipelineDirFactory,
    pathToPipelineFileFactory:pathToPipelineFileFactory,
    yamlCodec:yamlCodec,
  }

}

type _addPipelineUseCase struct {
  filesys                   ports.Filesys
  pathToPipelineDirFactory  pathToPipelineFileFactory
  pathToPipelineFileFactory pathToPipelineFileFactory
  yamlCodec                 yamlCodec
}

func (this _addPipelineUseCase) Execute(
req models.AddPipelineReq,
) (err error) {

  pathToPipelineDir := this.pathToPipelineDirFactory.Construct(
    req.ProjectUrl,
    req.Name,
  )

  err = this.filesys.CreateDir(pathToPipelineDir)
  if (nil != err) {
    return
  }

  var pipelineFile = pipelineFile{
    Description:req.Description,
  }

  pipelineFileBytes, err := this.yamlCodec.toYaml(&pipelineFile)
  if (nil != err) {
    return
  }

  pathToPipelineFile := this.pathToPipelineFileFactory.Construct(
    req.ProjectUrl,
    req.Name,
  )

  err = this.filesys.SaveFile(
    pathToPipelineFile,
    pipelineFileBytes,
  )

  return

}
