package osfilesys

import "io/ioutil"

type readDevOpFileUcExecuter interface {
  Execute(
  devOpName string,
  ) (devOpFile []byte, err error)
}

func newReadDevOpFileUcExecuter(
relPathToDevOpFileFactory relPathToDevOpDirFactory,
) readDevOpFileUcExecuter {

  return &readDevOpFileUcExecuterImpl{
    relPathToDevOpFileFactory:relPathToDevOpFileFactory,
  }

}

type readDevOpFileUcExecuterImpl struct {
  relPathToDevOpFileFactory relPathToDevOpDirFactory
}

func (uc readDevOpFileUcExecuterImpl) Execute(
devOpName string,
) (devOpFile []byte, err error) {

  relativePathToDevOpFile, err := uc.relPathToDevOpFileFactory.Construct(devOpName)
  if (nil != err) {
    return
  }

  devOpFile, err = ioutil.ReadFile(relativePathToDevOpFile)

  return
}
