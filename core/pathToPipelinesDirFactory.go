package core

import (
  "path"
  "github.com/dev-op-spec/engine/core/models"
)

type pathToPipelinesDirFactory interface {
  Construct(
  projectUrl *models.ProjectUrl,
  ) (pathToPipelinesDir string)
}

func newPathToPipelinesDirFactory(
) (pathToPipelinesDirFactory pathToPipelinesDirFactory) {

  return &_pathToPipelinesDirFactory{}

}

type _pathToPipelinesDirFactory struct{}

func (this _pathToPipelinesDirFactory) Construct(
projectUrl *models.ProjectUrl,
) (pathToPipelinesDir string) {

  pathToPipelinesDir = path.Join(
    projectUrl.Path,
    ".dev-op-spec",
    "pipelines",
  )

  return

}
