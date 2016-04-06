package core

import (
  "path"
  "github.com/dev-op-spec/engine/core/models"
)

type pathToOpDirFactory interface {
  Construct(
  projectUrl *models.Url,
  opName *string,
  ) (pathToOpDir string)
}

func newPathToOpDirFactory(
pathToOpsDirFactory pathToOpsDirFactory,
) (pathToOpDirFactory pathToOpFileFactory) {

  return &_pathToOpDirFactory{
    pathToOpsDirFactory:pathToOpsDirFactory,
  }

}

type _pathToOpDirFactory struct {
  pathToOpsDirFactory pathToOpsDirFactory
}

func (this _pathToOpDirFactory) Construct(
projectUrl *models.Url,
opName *string,
) (pathToOpDir string) {

  pathToOpsDir := this.pathToOpsDirFactory.Construct(
    projectUrl,
  )

  pathToOpDir = path.Join(
    pathToOpsDir,
    *opName,
  )

  return

}
