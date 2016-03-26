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
  return &_relPathToPipelineDirFactory{}
}

type _relPathToPipelineDirFactory struct{}

func (this _relPathToPipelineDirFactory) Construct(
pipelineName string,
) (relPathToPipelineDir string, err error) {

  if ("" == pipelineName) {
    err = errors.New("pipelineName cannot be nil")
  }

  relPathToPipelineDir = path.Join(relPathToPipelinesDir, pipelineName)

  return
}