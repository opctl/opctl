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

  return &_relPathToDevOpFileFactory{
    relPathToDevOpDirFactory:relPathToDevOpDirFactory,
  }

}

type _relPathToDevOpFileFactory struct {
  relPathToDevOpDirFactory relPathToDevOpDirFactory
}

func (this _relPathToDevOpFileFactory) Construct(
devOpName string,
) (relPathToDevOpFile string, err error) {

  var relPathToDevOpDir string
  relPathToDevOpDir, err = this.relPathToDevOpDirFactory.Construct(devOpName)
  if (nil != err) {
    return
  }

  relPathToDevOpFile = path.Join(relPathToDevOpDir, "dev-op.yml")

  return

}