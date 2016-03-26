package osfilesys

import "os"

type createDevOpDirUseCase interface {
  Execute(
  devOpName string,
  ) (err error)
}

func newCreateDevOpDirUseCase(
relPathToDevOpDirFactory relPathToDevOpDirFactory,
) createDevOpDirUseCase {

  return &_createDevOpDirUseCase{
    relPathToDevOpDirFactory:relPathToDevOpDirFactory,
  }

}

type _createDevOpDirUseCase struct {
  relPathToDevOpDirFactory relPathToDevOpDirFactory
}

func (this _createDevOpDirUseCase)  Execute(
devOpName string,
) (err error) {

  var relPathToDevOpDir string
  relPathToDevOpDir, err = this.relPathToDevOpDirFactory.Construct(devOpName)
  if (nil != err) {
    return
  }

  err = os.MkdirAll(relPathToDevOpDir, 0777)

  return

}