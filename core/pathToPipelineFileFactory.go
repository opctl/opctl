package core

import (
  "path"
  "github.com/dev-op-spec/engine/core/models"
)

type pathToPipelineFileFactory interface {
  Construct(
  projectUrl *models.ProjectUrl,
  pipelineName string,
  ) (pathToPipelineFile string)
}

func newPathToPipelineFileFactory(
pathToPipelineDirFactory pathToPipelineFileFactory,
) (pathToPipelineFileFactory pathToPipelineFileFactory) {

  return &_pathToPipelineFileFactory{
    pathToPipelineDirFactory:pathToPipelineDirFactory,
  }

}

type _pathToPipelineFileFactory struct {
  pathToPipelineDirFactory pathToPipelineFileFactory
}

func (this _pathToPipelineFileFactory) Construct(
projectUrl *models.ProjectUrl,
pipelineName string,
) (pathToPipelineFile string) {

  pathToPipelineDir := this.pathToPipelineDirFactory.Construct(
    projectUrl,
    pipelineName,
  )

  pathToPipelineFile = path.Join(
    pathToPipelineDir,
    "pipeline.yml",
  )

  return

}
