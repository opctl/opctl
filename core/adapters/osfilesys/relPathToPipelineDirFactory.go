package osfilesys

import (
  "path"
  "errors"
)

type relPathToPipelineDirFactory interface {
  Construct(
  pipelineName string,
  ) (relPathToPipelineDir string, err error)
}

func newRelPathToPipelineDirFactory() relPathToPipelineDirFactory {
  return &relPathToPipelineDirFactoryImpl{}
}

type relPathToPipelineDirFactoryImpl struct{}

func (f relPathToPipelineDirFactoryImpl) Construct(
pipelineName string,
) (relPathToPipelineDir string, err error) {

  if ("" == pipelineName) {
    err = errors.New("pipelineName cannot be nil")
  }

  relPathToPipelineDir = path.Join(relPathToPipelinesDir, pipelineName)

  return
}