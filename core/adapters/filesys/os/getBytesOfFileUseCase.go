package os

import "io/ioutil"

type getBytesOfFileUseCase interface {
  Execute(
  pathToFile string,
  ) (bytesOfFile []byte, err error)
}

func newGetBytesOfFileUseCase(
) getBytesOfFileUseCase {

  return &_getBytesOfFileUseCase{}

}

type _getBytesOfFileUseCase struct{}

func (this _getBytesOfFileUseCase) Execute(
pathToFile string,
) (bytesOfFile []byte, err error) {

  bytesOfFile, err = ioutil.ReadFile(pathToFile)

  return
}
