package core

import (
  "path"
  "github.com/dev-op-spec/engine/core/models"
)

type pathToPipelineDirFactory interface {
  Construct(
  projectUrl *models.ProjectUrl,
  pipelineName string,
  ) (pathToPipelineDir string)
}

func newPathToPipelineDirFactory(
pathToPipelinesDirFactory pathToPipelinesDirFactory,
) (pathToPipelineDirFactory pathToPipelineFileFactory) {

  return &_pathToPipelineDirFactory{
    pathToPipelinesDirFactory:pathToPipelinesDirFactory,
  }

}

type _pathToPipelineDirFactory struct {
  pathToPipelinesDirFactory pathToPipelinesDirFactory
}

func (this _pathToPipelineDirFactory) Construct(
projectUrl *models.ProjectUrl,
pipelineName string,
) (pathToPipelineDir string) {

  pathToPipelinesDir := this.pathToPipelinesDirFactory.Construct(
    projectUrl,
  )

  pathToPipelineDir = path.Join(
    pathToPipelinesDir,
    pipelineName,
  )

  return

}
