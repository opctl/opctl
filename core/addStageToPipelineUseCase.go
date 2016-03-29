package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type addStageToPipelineUseCase interface {
  Execute(
  req models.AddStageToPipelineReq,
  ) (err error)
}

func newAddStageToPipelineUseCase(
filesys ports.Filesys,
pathToPipelineFileFactory pathToPipelineFileFactory,
yamlCodec yamlCodec,
) addStageToPipelineUseCase {

  return &_addStageToPipelineUseCase{
    filesys:filesys,
    pathToPipelineFileFactory:pathToPipelineFileFactory,
    yamlCodec:yamlCodec,
  }

}

type _addStageToPipelineUseCase struct {
  filesys                   ports.Filesys
  pathToPipelineFileFactory pathToPipelineFileFactory
  yamlCodec                 yamlCodec
}

func (this _addStageToPipelineUseCase) Execute(
req models.AddStageToPipelineReq,
) (err error) {

  pathToPipelineFile := this.pathToPipelineFileFactory.Construct(
    req.PathToProjectRootDir,
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

  newPipelineFileStage := pipelineFileStage{Name:req.StageName}

  // set type
  if (req.IsPipelineStage) {
    newPipelineFileStage.Type = pipelineStageType
  } else {
    newPipelineFileStage.Type = devOpStageType
  }

  if (nil == pipelineFile.Stages) {
    pipelineFile.Stages = []pipelineFileStage{}
  }

  if (len(req.PrecedingStageName) > 0) {

    stages := []pipelineFileStage{}

    for _, stage := range pipelineFile.Stages {

      stages = append(stages, stage)
      if (stage.Name == req.PrecedingStageName) {
        stages = append(stages, newPipelineFileStage)
      }

    }

    pipelineFile.Stages = stages

  } else {

    pipelineFile.Stages = append(pipelineFile.Stages, newPipelineFileStage)

  }

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