package core

import "path"

type pathToDevOpDirFactory interface {
  Construct(
  pathToProjectRootDir string,
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
pathToProjectRootDir string,
devOpName string,
) (pathToDevOpDir string) {

  pathToDevOpsDir := this.pathToDevOpsDirFactory.Construct(
    pathToProjectRootDir,
  )

  pathToDevOpDir = path.Join(
    pathToDevOpsDir,
    devOpName,
  )

  return

}
