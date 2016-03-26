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
fs ports.Filesys,
yml yamlCodec,
) addStageToPipelineUseCase {

  return &_addStageToPipelineUseCase{
    fs:fs,
    yml:yml,
  }

}

type _addStageToPipelineUseCase struct {
  fs  ports.Filesys
  yml yamlCodec
}

func (this _addStageToPipelineUseCase) Execute(
req models.AddStageToPipelineReq,
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

  newPipelineFileStage := pipelineFileStage{Name:req.StageName}

  // set type
  if (req.IsPipelineStage) {
    newPipelineFileStage.Type = pipelineStageType
  } else {
    newPipelineFileStage.Type = devOpStageType
  }

  if (nil == pipelineFile.Stages) {
    pipelineFile.Stages = make([]pipelineFileStage, 0)
  }

  if (len(req.PrecedingStageName) > 0) {

    var stages = make([]pipelineFileStage, 0)

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