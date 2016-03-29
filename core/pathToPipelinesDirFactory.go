package core

import "path"

type pathToPipelinesDirFactory interface {
  Construct(
  pathToProjectRootDir string,
  ) (pathToPipelinesDir string)
}

func newPathToPipelinesDirFactory(
) (pathToPipelinesDirFactory pathToPipelinesDirFactory) {

  return &_pathToPipelinesDirFactory{}

}

type _pathToPipelinesDirFactory struct{}

func (this _pathToPipelinesDirFactory) Construct(
pathToProjectRootDir string,
) (pathToPipelinesDir string) {

  pathToPipelinesDir = path.Join(
    pathToProjectRootDir,
    ".dev-op-spec",
    "pipelines",
  )

  return

}
