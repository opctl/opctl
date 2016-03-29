package core

import (
  "path"
)

type pathToDevOpFileFactory interface {
  Construct(
  pathToProjectRootDir string,
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
pathToProjectRootDir string,
devOpName string,
) (pathToDevOpFile string) {

  pathToDevOpDir := this.pathToDevOpDirFactory.Construct(
    pathToProjectRootDir,
    devOpName,
  )

  pathToDevOpFile = path.Join(
    pathToDevOpDir,
    "dev-op.yml",
  )

  return

}
