package osfilesys

import "io/ioutil"

type saveDevOpFileUcExecuter interface {
  Execute(
  devOpName string,
  data []byte,
  ) (err error)
}

func newSaveDevOpFileUcExecuter(
relPathToDevOpFileFactory relPathToDevOpFileFactory,
) saveDevOpFileUcExecuter {

  return &saveDevOpFileUcExecuterImpl{
    relPathToDevOpFileFactory:relPathToDevOpFileFactory,
  }

}

type saveDevOpFileUcExecuterImpl struct {
  relPathToDevOpFileFactory relPathToDevOpFileFactory
}

func (uc saveDevOpFileUcExecuterImpl)  Execute(
devOpName string,
data []byte,
) (err error) {

  var relPathToDevOpFile string
  relPathToDevOpFile, err = uc.relPathToDevOpFileFactory.Construct(devOpName)
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