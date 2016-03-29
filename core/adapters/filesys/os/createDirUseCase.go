package os

import "os"

type createDirUseCase interface {
  Execute(
  devOpName string,
  ) (err error)
}

func newCreateDirUseCase(
) createDirUseCase {

  return &_createDirUseCase{}

}

type _createDirUseCase struct {
}

func (this _createDirUseCase)  Execute(
pathToDir string,
) (err error) {

  err = os.MkdirAll(pathToDir, 0777)

  return

}