package core

import (
  "path"
  "github.com/dev-op-spec/engine/core/models"
)

type pathToDevOpDirFactory interface {
  Construct(
  projectUrl *models.ProjectUrl,
  devOpName string,
  ) (pathToDevOpDir string)
}

func newPathToDevOpDirFactory(
pathToDevOpsDirFactory pathToDevOpsDirFactory,
) (pathToDevOpDirFactory pathToDevOpDirFactory) {

  return &_pathToDevOpDirFactory{
    pathToDevOpsDirFactory:pathToDevOpsDirFactory,
  }

}

type _pathToDevOpDirFactory struct {
  pathToDevOpsDirFactory pathToDevOpsDirFactory
}

func (this _pathToDevOpDirFactory) Construct(
projectUrl *models.ProjectUrl,
devOpName string,
) (pathToDevOpDir string) {

  pathToDevOpsDir := this.pathToDevOpsDirFactory.Construct(
    projectUrl,
  )

  pathToDevOpDir = path.Join(
    pathToDevOpsDir,
    devOpName,
  )

  return

}
