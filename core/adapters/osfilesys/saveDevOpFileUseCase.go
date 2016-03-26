package osfilesys

import "io/ioutil"

type saveDevOpFileUseCase interface {
  Execute(
  devOpName string,
  data []byte,
  ) (err error)
}

func newSaveDevOpFileUseCase(
relPathToDevOpFileFactory relPathToDevOpFileFactory,
) saveDevOpFileUseCase {

  return &_saveDevOpFileUseCase{
    relPathToDevOpFileFactory:relPathToDevOpFileFactory,
  }

}

type _saveDevOpFileUseCase struct {
  relPathToDevOpFileFactory relPathToDevOpFileFactory
}

func (this _saveDevOpFileUseCase)  Execute(
devOpName string,
data []byte,
) (err error) {

  var relPathToDevOpFile string
  relPathToDevOpFile, err = this.relPathToDevOpFileFactory.Construct(devOpName)
  if (nil != err) {
    return
  }

  err = ioutil.WriteFile(
    relPathToDevOpFile,
    data,
    0777,
  )

  return

}