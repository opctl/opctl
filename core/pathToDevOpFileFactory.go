package core

import (
  "path"
  "github.com/dev-op-spec/engine/core/models"
)

type pathToDevOpFileFactory interface {
  Construct(
  projectUrl *models.ProjectUrl,
  devOpName string,
  ) (pathToDevOpFile string)
}

func newPathToDevOpFileFactory(
pathToDevOpDirFactory pathToDevOpFileFactory,
) (pathToDevOpFileFactory pathToDevOpFileFactory) {

  return &_pathToDevOpFileFactory{
    pathToDevOpDirFactory:pathToDevOpDirFactory,
  }

}

type _pathToDevOpFileFactory struct {
  pathToDevOpDirFactory pathToDevOpFileFactory
}

func (this _pathToDevOpFileFactory) Construct(
projectUrl *models.ProjectUrl,
devOpName string,
) (pathToDevOpFile string) {

  pathToDevOpDir := this.pathToDevOpDirFactory.Construct(
    projectUrl,
    devOpName,
  )

  pathToDevOpFile = path.Join(
    pathToDevOpDir,
    "dev-op.yml",
  )

  return

}
