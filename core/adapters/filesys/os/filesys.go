package os

import "github.com/open-devops/engine/core/ports"

func New() ports.Filesys {

  return filesys{
    compositionRoot:newCompositionRoot(),
  }

}

type filesys struct {
  compositionRoot compositionRoot
}

func (this filesys) CreateDir(
pathToDir string,
) (err error) {
  return this.compositionRoot.
  CreateDirUseCase().
  Execute(pathToDir)
}

func (this filesys) ListNamesOfChildDirs(
pathToParentDir string,
) (namesOfChildDirs []string, err error) {
  return this.compositionRoot.
  ListNamesOfChildDirsUseCase().
  Execute(pathToParentDir)
}

func (this filesys) GetBytesOfFile(
pathToFile string,
) (bytesOfFile []byte, err error) {
  return this.compositionRoot.
  GetBytesOfFileUseCase().
  Execute(pathToFile)
}

func (this filesys) SaveFile(
pathToFile string,
bytesOfFile []byte,
) (err error) {
  return this.compositionRoot.
  SaveFileUseCase().
  Execute(
    pathToFile,
    bytesOfFile,
  )
}
