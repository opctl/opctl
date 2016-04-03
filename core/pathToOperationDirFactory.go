package core

import (
  "path"
  "github.com/dev-op-spec/engine/core/models"
)

type pathToOperationDirFactory interface {
  Construct(
  projectUrl *models.ProjectUrl,
  operationName string,
  ) (pathToOperationDir string)
}

func newPathToOperationDirFactory(
pathToOperationsDirFactory pathToOperationsDirFactory,
) (pathToOperationDirFactory pathToOperationFileFactory) {

  return &_pathToOperationDirFactory{
    pathToOperationsDirFactory:pathToOperationsDirFactory,
  }

}

type _pathToOperationDirFactory struct {
  pathToOperationsDirFactory pathToOperationsDirFactory
}

func (this _pathToOperationDirFactory) Construct(
projectUrl *models.ProjectUrl,
operationName string,
) (pathToOperationDir string) {

  pathToOperationsDir := this.pathToOperationsDirFactory.Construct(
    projectUrl,
  )

  pathToOperationDir = path.Join(
    pathToOperationsDir,
    operationName,
  )

  return

}
