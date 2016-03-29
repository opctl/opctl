package os

import "io/ioutil"

type saveFileUseCase interface {
  Execute(
  pathOfFile string,
  bytesOfFile []byte,
  ) (err error)
}

func newSaveFileUseCase(
) saveFileUseCase {

  return &_saveFileUseCase{}

}

type _saveFileUseCase struct{}

func (this _saveFileUseCase)  Execute(
pathOfFile string,
bytesOfFile []byte,
) (err error) {

  err = ioutil.WriteFile(
    pathOfFile,
    bytesOfFile,
    0777,
  )

  return

}