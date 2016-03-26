package os

import (
  "strings"
  "io/ioutil"
  "os"
)

type listNamesOfPipelineDirsUseCase interface {
  Execute() (namesOfPipelineDirs []string, err error)
}

func newListNamesOfPipelineDirsUseCase() listNamesOfPipelineDirsUseCase {

  return &_listNamesOfPipelineDirsUseCase{}

}

type _listNamesOfPipelineDirsUseCase struct{}

func (this _listNamesOfPipelineDirsUseCase)  Execute(
) (namesOfPipelineDirs []string, err error) {

  namesOfPipelineDirs = make([]string, 0)

  var pipelineDirFileInfos []os.FileInfo
  pipelineDirFileInfos, err = ioutil.ReadDir(relPathToPipelinesDir)
  if (nil != err) {
    return
  }

  for _, x := range pipelineDirFileInfos {

    if x.IsDir() && !strings.HasPrefix(x.Name(), ".") {

      namesOfPipelineDirs = append(namesOfPipelineDirs, x.Name())

    }

  }

  return

}
