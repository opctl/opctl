package core

import (
  "path"
  "github.com/dev-op-spec/engine/core/models"
)

type pathToOpFileFactory interface {
  Construct(
  projectUrl *models.Url,
  opName *string,
  ) (pathToOpFile string)
}

func newPathToOpFileFactory(
pathToOpDirFactory pathToOpFileFactory,
) (pathToOpFileFactory pathToOpFileFactory) {

  return &_pathToOpFileFactory{
    pathToOpDirFactory:pathToOpDirFactory,
  }

}

type _pathToOpFileFactory struct {
  pathToOpDirFactory pathToOpFileFactory
}

func (this _pathToOpFileFactory) Construct(
projectUrl *models.Url,
opName *string,
) (pathToOpFile string) {

  pathToOpDir := this.pathToOpDirFactory.Construct(
    projectUrl,
    opName,
  )

  pathToOpFile = path.Join(
    pathToOpDir,
    "op.yml",
  )

  return

}
