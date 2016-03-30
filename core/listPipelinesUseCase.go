package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type listPipelinesUseCase interface {
  Execute(
  projectUrl *models.ProjectUrl,
  ) (pipelines []models.PipelineView, err error)
}

func newListPipelinesUseCase(
filesys ports.Filesys,
pathToPipelineFileFactory pathToPipelineFileFactory,
pathToPipelinesDirFactory pathToPipelinesDirFactory,
yamlCodec yamlCodec,
) listPipelinesUseCase {

  return &_listPipelinesUseCase{
    filesys:filesys,
    pathToPipelineFileFactory:pathToPipelineFileFactory,
    pathToPipelinesDirFactory:pathToPipelinesDirFactory,
    yamlCodec:yamlCodec,
  }

}

type _listPipelinesUseCase struct {
  filesys                   ports.Filesys
  pathToPipelineFileFactory pathToPipelineFileFactory
  pathToPipelinesDirFactory pathToPipelinesDirFactory
  yamlCodec                 yamlCodec
}

func (this _listPipelinesUseCase) Execute(
projectUrl *models.ProjectUrl,
) (pipelines []models.PipelineView, err error) {

  pathToPipelinesDir := this.pathToPipelinesDirFactory.Construct(
    projectUrl,
  )

  pipelineDirNames, err := this.filesys.ListNamesOfChildDirs(
    pathToPipelinesDir,
  )
  if (nil != err) {
    return
  }

  for _, pipelineDirName := range pipelineDirNames {

    pathToPipelineFile := this.pathToPipelineFileFactory.Construct(
      projectUrl,
      pipelineDirName,
    )

    var pipelineFileBytes []byte
    pipelineFileBytes, err = this.filesys.GetBytesOfFile(pathToPipelineFile)
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

    pipelineStageViews := []models.PipelineStageView{}

    for _, pipelineStage := range pipelineFile.Stages {

      pipelineStageView := models.NewPipelineStageView(
        pipelineStage.Name,
        pipelineStage.Type,
      )

      pipelineStageViews = append(pipelineStageViews, *pipelineStageView)

    }

    pipelineView := models.NewPipelineView(
      pipelineFile.Description,
      pipelineDirName,
      pipelineStageViews,
    )

    pipelines = append(pipelines, *pipelineView)

  }

  return

}
