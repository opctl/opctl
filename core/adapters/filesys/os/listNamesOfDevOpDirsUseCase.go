package os

import (
  "strings"
  "io/ioutil"
  "os"
)

type listNamesOfDevOpDirsUseCase interface {
  Execute() (namesOfDevOpDirs []string, err error)
}

func newListNamesOfDevOpDirsUseCase() listNamesOfDevOpDirsUseCase {

  return &_listNamesOfDevOpDirsUseCase{}

}

type _listNamesOfDevOpDirsUseCase struct{}

func (this _listNamesOfDevOpDirsUseCase)  Execute(
) (namesOfDevOpDirs []string, err error) {

  namesOfDevOpDirs = make([]string, 0)

  var devOpDirFileInfos []os.FileInfo
  devOpDirFileInfos, err = ioutil.ReadDir(relPathToDevOpsDir)
  if (nil != err) {
    return
  }

  for _, x := range devOpDirFileInfos {

    if x.IsDir() && !strings.HasPrefix(x.Name(), ".") {

      namesOfDevOpDirs = append(namesOfDevOpDirs, x.Name())

    }

  }

  return

}
