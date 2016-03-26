package os

import "io/ioutil"

type readDevOpFileUseCase interface {
  Execute(
  devOpName string,
  ) (devOpFile []byte, err error)
}

func newReadDevOpFileUseCase(
relPathToDevOpFileFactory relPathToDevOpDirFactory,
) readDevOpFileUseCase {

  return &_readDevOpFileUseCase{
    relPathToDevOpFileFactory:relPathToDevOpFileFactory,
  }

}

type _readDevOpFileUseCase struct {
  relPathToDevOpFileFactory relPathToDevOpDirFactory
}

func (this _readDevOpFileUseCase) Execute(
devOpName string,
) (devOpFile []byte, err error) {

  relativePathToDevOpFile, err := this.relPathToDevOpFileFactory.Construct(devOpName)
  if (nil != err) {
    return
  }

  devOpFile, err = ioutil.ReadFile(relativePathToDevOpFile)

  return
}
