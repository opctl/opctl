package core

import "path"

type pathToDevOpsDirFactory interface {
  Construct(
  pathToProjectRootDir string,
  ) (pathToDevOpsDir string)
}

func newPathToDevOpsDirFactory(
) (pathToDevOpsDirFactory pathToDevOpsDirFactory) {

  return &_pathToDevOpsDirFactory{}

}

type _pathToDevOpsDirFactory struct{}

func (this _pathToDevOpsDirFactory) Construct(
pathToProjectRootDir string,
) (pathToDevOpsDir string) {

  pathToDevOpsDir = path.Join(
    pathToProjectRootDir,
    ".dev-op-spec",
    "dev-ops",
  )

  return

}
