package core

import (
  "path"
  "github.com/dev-op-spec/engine/core/models"
)

type pathToOperationFileFactory interface {
  Construct(
  projectUrl *models.Url,
  operationName *string,
  ) (pathToOperationFile string)
}

func newPathToOperationFileFactory(
pathToOperationDirFactory pathToOperationFileFactory,
) (pathToOperationFileFactory pathToOperationFileFactory) {

  return &_pathToOperationFileFactory{
    pathToOperationDirFactory:pathToOperationDirFactory,
  }

}

type _pathToOperationFileFactory struct {
  pathToOperationDirFactory pathToOperationFileFactory
}

func (this _pathToOperationFileFactory) Construct(
projectUrl *models.Url,
operationName *string,
) (pathToOperationFile string) {

  pathToOperationDir := this.pathToOperationDirFactory.Construct(
    projectUrl,
    operationName,
  )

  pathToOperationFile = path.Join(
    pathToOperationDir,
    "operation.yml",
  )

  return

}
