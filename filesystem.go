package opspec

//go:generate counterfeiter -o ./fakeFileSystem.go --fake-name FakeFilesystem ./ Filesystem

import (
  "io/ioutil"
  "os"
)

type Filesystem  interface {
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

func newFilesystem(
) Filesystem {
  return &_filesystem{}
}

type _filesystem struct{}

func (this _filesystem)  AddDir(
pathToDir string,
) (err error) {

  err = os.MkdirAll(pathToDir, 0777)

  return

}

func (this _filesystem) GetBytesOfFile(
pathToFile string,
) (bytesOfFile []byte, err error) {

  bytesOfFile, err = ioutil.ReadFile(pathToFile)

  return
}

func (this _filesystem) ListChildFileInfosOfDir(
pathToDir string,
) (
childFileInfos []os.FileInfo,
err error,
) {
  childFileInfos, err = ioutil.ReadDir(pathToDir)
  return
}

func (this _filesystem) SaveFile(
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
