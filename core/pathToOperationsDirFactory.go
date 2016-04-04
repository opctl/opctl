package core

import (
  "path"
  "github.com/dev-op-spec/engine/core/models"
)

type pathToOperationsDirFactory interface {
  Construct(
  projectUrl *models.Url,
  ) (pathToOperationsDir string)
}

func newPathToOperationsDirFactory(
) (pathToOperationsDirFactory pathToOperationsDirFactory) {

  return &_pathToOperationsDirFactory{}

}

type _pathToOperationsDirFactory struct{}

func (this _pathToOperationsDirFactory) Construct(
projectUrl *models.Url,
) (pathToOperationsDir string) {

  pathToOperationsDir = path.Join(
    projectUrl.Path,
    ".dev-op-spec",
    "operations",
  )

  return

}
