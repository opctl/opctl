package core

import (
  "path"
)

type pathToPipelineFileFactory interface {
  Construct(
  pathToProjectRootDir string,
  pipelineName string,
  ) (pathToPipelineFile string)
}

func newPathToPipelineFileFactory(
pathToPipelineDirFactory pathToPipelineFileFactory,
) (pathToPipelineFileFactory pathToPipelineFileFactory) {

  return &_pathToPipelineFileFactory{
    pathToPipelineDirFactory:pathToPipelineDirFactory,
  }

}

type _pathToPipelineFileFactory struct {
  pathToPipelineDirFactory pathToPipelineFileFactory
}

func (this _pathToPipelineFileFactory) Construct(
pathToProjectRootDir string,
pipelineName string,
) (pathToPipelineFile string) {

  pathToPipelineDir := this.pathToPipelineDirFactory.Construct(
    pathToProjectRootDir,
    pipelineName,
  )

  pathToPipelineFile = path.Join(
    pathToPipelineDir,
    "pipeline.yml",
  )

  return

}
