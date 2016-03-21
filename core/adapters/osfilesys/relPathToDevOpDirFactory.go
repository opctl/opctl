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
  return &relPathToDevOpDirFactoryImpl{}
}

type relPathToDevOpDirFactoryImpl struct{}

func (f relPathToDevOpDirFactoryImpl) Construct(
devOpName string,
) (relPathToDevOpDir string, err error) {

  if ("" == devOpName) {
    err = errors.New("devOpName cannot be nil")
  }

  relPathToDevOpDir = path.Join(relPathToDevOpsDir, devOpName)

  return
}