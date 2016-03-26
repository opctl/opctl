package osfilesys

import (
  "path"
  "errors"
)

type relPathToDevOpDirFactory interface {
  Construct(
  devOpName string,
  ) (relPathToDevOpDir string, err error)
}

func newRelPathToDevOpDirFactory() relPathToDevOpDirFactory {
  return &_relPathToDevOpDirFactory{}
}

type _relPathToDevOpDirFactory struct{}

func (this _relPathToDevOpDirFactory) Construct(
devOpName string,
) (relPathToDevOpDir string, err error) {

  if ("" == devOpName) {
    err = errors.New("devOpName cannot be nil")
  }

  relPathToDevOpDir = path.Join(relPathToDevOpsDir, devOpName)

  return
}