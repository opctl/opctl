package os

import (
  "strings"
  "io/ioutil"
)

type listNamesOfChildDirsUseCase interface {
  Execute(
  pathToParentDir string,
  ) (namesOfChildDirs []string, err error)
}

func newListNamesOfChildDirsUseCase() listNamesOfChildDirsUseCase {

  return &_listNamesOfChildDirsUseCase{}

}

type _listNamesOfChildDirsUseCase struct{}

func (this _listNamesOfChildDirsUseCase)  Execute(
pathToParentDir string,
) (namesOfChildDirs []string, err error) {

  namesOfChildDirs = []string{}

  childDirFileInfos, err := ioutil.ReadDir(pathToParentDir)
  if (nil != err) {
    return
  }

  for _, x := range childDirFileInfos {

    if x.IsDir() && !strings.HasPrefix(x.Name(), ".") {

      namesOfChildDirs = append(namesOfChildDirs, x.Name())

    }

  }

  return

}
