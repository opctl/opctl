package fs

//go:generate counterfeiter -o ./fakeFileSystem.go --fake-name FakeFileSystem ./ FileSystem

import (
  "io/ioutil"
  "os"
)

type FileSystem  interface {
  AddDir(
  pathToDir string,
  ) (err error)

  GetBytesOfFile(
  pathToFile string,
  ) (
  bytesOfFile []byte,
  err error,
  )

  ListChildFileInfosOfDir(
  pathToDir string,
  ) (
  childFileInfos []os.FileInfo,
  err error,
  )

  SaveFile(
  pathToFile string,
  bytesOfFile []byte,
  ) (err error)
}

func NewFileSystem(
) FileSystem {
  return &_fileSystem{}
}

type _fileSystem struct{}

func (this _fileSystem)  AddDir(
pathToDir string,
) (err error) {

  err = os.MkdirAll(pathToDir, 0777)

  return

}

func (this _fileSystem) GetBytesOfFile(
pathToFile string,
) (bytesOfFile []byte, err error) {

  bytesOfFile, err = ioutil.ReadFile(pathToFile)

  return
}

func (this _fileSystem) ListChildFileInfosOfDir(
pathToDir string,
) (
childFileInfos []os.FileInfo,
err error,
) {
  childFileInfos, err = ioutil.ReadDir(pathToDir)
  return
}

func (this _fileSystem) SaveFile(
pathToFile string,
bytesOfFile []byte,
) (err error) {

  err = ioutil.WriteFile(
    pathToFile,
    bytesOfFile,
    0777,
  )

  return
}
