package core

import (
  "github.com/dev-op-spec/engine/core/models"
"github.com/dev-op-spec/engine/core/ports"
)

type listDevOpsUcExecuter interface {
  Execute(
  ) (devOps []models.DevOpView, err error)
}

func newListDevOpsUcExecuter(
fs ports.Filesys,
yml yamlCodec,
) listDevOpsUcExecuter {

  return &listDevOpsUcExecuterImpl{
    fs:fs,
    yml:yml,
  }

}

type listDevOpsUcExecuterImpl struct {
  fs  ports.Filesys
  yml yamlCodec
}

func (x listDevOpsUcExecuterImpl) Execute(
) (devOps []models.DevOpView, err error) {

  devOps = make([]models.DevOpView, 0)

  var devOpDirNames []string
  devOpDirNames, err= x.fs.ListNamesOfDevOpDirs()
  if (nil != err) {
    return
  }

  for _, devOpDirName := range devOpDirNames {

    var devOpFileBytes []byte
    devOpFileBytes, err= x.fs.ReadDevOpFile(devOpDirName)
    if (nil != err) {
      return
    }

    devOpFile := devOpFile{}
    err= x.yml.fromYaml(
      devOpFileBytes,
      &devOpFile,
    )
    if (nil != err) {
      return
    }

    devOpView := models.NewDevOpView(devOpFile.Description, devOpDirName)

    devOps = append(devOps, *devOpView)

  }

  return

}
