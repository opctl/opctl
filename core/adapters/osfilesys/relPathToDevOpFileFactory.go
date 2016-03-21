package osfilesys

import "path"

type relPathToDevOpFileFactory interface {
  Construct(
  devOpName string,
  ) (relPathToDevOpFile string, err error)
}

func newRelPathToDevOpFileFactory(
relPathToDevOpDirFactory relPathToDevOpDirFactory,
) relPathToDevOpFileFactory {

  return &relPathToDevOpFileFactoryImpl{
    relPathToDevOpDirFactory:relPathToDevOpDirFactory,
  }

}

type relPathToDevOpFileFactoryImpl struct {
  relPathToDevOpDirFactory relPathToDevOpDirFactory
}

func (f relPathToDevOpFileFactoryImpl) Construct(
devOpName string,
) (relPathToDevOpFile string, err error) {

  var relPathToDevOpDir string
  relPathToDevOpDir, err = f.relPathToDevOpDirFactory.Construct(devOpName)
  if (nil != err) {
    return
  }

  relPathToDevOpFile = path.Join(relPathToDevOpDir, "dev-op.yml")

  return

}