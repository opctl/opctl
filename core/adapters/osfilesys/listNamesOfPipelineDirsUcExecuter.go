package osfilesys

import (
  "strings"
  "io/ioutil"
  "os"
)

type listNamesOfPipelineDirsUcExecuter interface {
  Execute() (namesOfPipelineDirs []string, err error)
}

func newListNamesOfPipelineDirsUcExecuter() listNamesOfPipelineDirsUcExecuter {

  return &listNamesOfPipelineDirsUcExecuterImpl{}

}

type listNamesOfPipelineDirsUcExecuterImpl struct{}

func (uc listNamesOfPipelineDirsUcExecuterImpl)  Execute(
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
