package core

import (
  "github.com/dev-op-spec/engine/core/models"
"github.com/dev-op-spec/engine/core/ports"
)

type listPipelinesUseCase interface {
  Execute(
  ) (pipelines []models.PipelineView, err error)
}

func newListPipelinesUseCase(
fs ports.Filesys,
yml yamlCodec,
) listPipelinesUseCase {

  return &_listPipelinesUseCase{
    fs:fs,
    yml:yml,
  }

}

type _listPipelinesUseCase struct {
  fs  ports.Filesys
  yml yamlCodec
}

func (this _listPipelinesUseCase) Execute(
) (pipelines []models.PipelineView, err error) {

  pipelines = make([]models.PipelineView, 0)

  var pipelineDirNames []string
  pipelineDirNames, err= this.fs.ListNamesOfPipelineDirs()
  if (nil != err) {
    return
  }

  for _, pipelineDirName := range pipelineDirNames {

    var pipelineFileBytes []byte
    pipelineFileBytes, err= this.fs.ReadPipelineFile(pipelineDirName)
    if (nil != err) {
      return
    }

    pipelineFile := pipelineFile{}
    err= this.yml.fromYaml(
      pipelineFileBytes,
      &pipelineFile,
    )
    if (nil != err) {
      return
    }

    pipelineStageViews := []models.PipelineStageView{}

    for _, pipelineStage := range pipelineFile.Stages {
      pipelineStageView := models.NewPipelineStageView(pipelineStage.Name, pipelineStage.Type)
      pipelineStageViews = append(pipelineStageViews, *pipelineStageView)
    }

    pipelineView := models.NewPipelineView(pipelineFile.Description, pipelineDirName, pipelineStageViews)

    pipelines = append(pipelines, *pipelineView)

  }

  return

}
