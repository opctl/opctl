package osfilesys

import (
  "strings"
  "io/ioutil"
  "os"
)

type listNamesOfDevOpDirsUcExecuter interface {
  Execute() (namesOfDevOpDirs []string, err error)
}

func newListNamesOfDevOpDirsUcExecuter() listNamesOfDevOpDirsUcExecuter {

  return &listNamesOfDevOpDirsUcExecuterImpl{}

}

type listNamesOfDevOpDirsUcExecuterImpl struct{}

func (uc listNamesOfDevOpDirsUcExecuterImpl)  Execute(
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
