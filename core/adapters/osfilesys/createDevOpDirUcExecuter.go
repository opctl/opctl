package osfilesys

import "os"

type createDevOpDirUcExecuter interface {
  Execute(
  devOpName string,
  ) (err error)
}

func newCreateDevOpDirUcExecuter(
relPathToDevOpDirFactory relPathToDevOpDirFactory,
) createDevOpDirUcExecuter {

  return &createDevOpDirUcExecuterImpl{
    relPathToDevOpDirFactory:relPathToDevOpDirFactory,
  }

}

type createDevOpDirUcExecuterImpl struct {
  relPathToDevOpDirFactory relPathToDevOpDirFactory
}

func (uc createDevOpDirUcExecuterImpl)  Execute(
devOpName string,
) (err error) {

  var relPathToDevOpDir string
  relPathToDevOpDir, err = uc.relPathToDevOpDirFactory.Construct(devOpName)
  if (nil != err) {
    return
  }

  err = os.MkdirAll(relPathToDevOpDir, 0777)

  return

}