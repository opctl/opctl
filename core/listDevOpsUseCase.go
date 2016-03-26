package core

import (
  "github.com/dev-op-spec/engine/core/models"
"github.com/dev-op-spec/engine/core/ports"
)

type listDevOpsUseCase interface {
  Execute(
  ) (devOps []models.DevOpView, err error)
}

func newListDevOpsUseCase(
fs ports.Filesys,
yml yamlCodec,
) listDevOpsUseCase {

  return &_listDevOpsUseCase{
    fs:fs,
    yml:yml,
  }

}

type _listDevOpsUseCase struct {
  fs  ports.Filesys
  yml yamlCodec
}

func (this _listDevOpsUseCase) Execute(
) (devOps []models.DevOpView, err error) {

  devOps = make([]models.DevOpView, 0)

  var devOpDirNames []string
  devOpDirNames, err= this.fs.ListNamesOfDevOpDirs()
  if (nil != err) {
    return
  }

  for _, devOpDirName := range devOpDirNames {

    var devOpFileBytes []byte
    devOpFileBytes, err= this.fs.ReadDevOpFile(devOpDirName)
    if (nil != err) {
      return
    }

    devOpFile := devOpFile{}
    err= this.yml.fromYaml(
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
