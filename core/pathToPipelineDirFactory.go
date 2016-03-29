package core

import "path"

type pathToPipelineDirFactory interface {
  Construct(
  pathToProjectRootDir string,
  pipelineName string,
  ) (pathToPipelineDir string)
}

func newPathToPipelineDirFactory(
pathToPipelinesDirFactory pathToPipelinesDirFactory,
) (pathToPipelineDirFactory pathToPipelineFileFactory) {

  return &_pathToPipelineDirFactory{
    pathToPipelinesDirFactory:pathToPipelinesDirFactory,
  }

}

type _pathToPipelineDirFactory struct {
  pathToPipelinesDirFactory pathToPipelinesDirFactory
}

func (this _pathToPipelineDirFactory) Construct(
pathToProjectRootDir string,
pipelineName string,
) (pathToPipelineDir string) {

  pathToPipelinesDir := this.pathToPipelinesDirFactory.Construct(
    pathToProjectRootDir,
  )

  pathToPipelineDir = path.Join(
    pathToPipelinesDir,
    pipelineName,
  )

  return

}
