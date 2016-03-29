package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type listDevOpsUseCase interface {
  Execute(
  pathToProjectRoot string,
  ) (devOps []models.DevOpView, err error)
}

func newListDevOpsUseCase(
filesys ports.Filesys,
pathToDevOpFileFactory pathToDevOpFileFactory,
pathToDevOpsDirFactory pathToDevOpsDirFactory,
yamlCodec yamlCodec,
) listDevOpsUseCase {

  return &_listDevOpsUseCase{
    filesys:filesys,
    pathToDevOpFileFactory:pathToDevOpFileFactory,
    pathToDevOpsDirFactory:pathToDevOpsDirFactory,
    yamlCodec:yamlCodec,
  }

}

type _listDevOpsUseCase struct {
  filesys                ports.Filesys
  pathToDevOpFileFactory pathToDevOpFileFactory
  pathToDevOpsDirFactory pathToDevOpsDirFactory
  yamlCodec              yamlCodec
}

func (this _listDevOpsUseCase) Execute(
pathToProjectRoot string,
) (devOps []models.DevOpView, err error) {

  pathToDevOpsDir := this.pathToDevOpsDirFactory.Construct(
    pathToProjectRoot,
  )

  devOpDirNames, err := this.filesys.ListNamesOfChildDirs(
    pathToDevOpsDir,
  )
  if (nil != err) {
    return
  }

  for _, devOpDirName := range devOpDirNames {

    pathToDevOpFile := this.pathToDevOpFileFactory.Construct(
      pathToProjectRoot,
      devOpDirName,
    )

    var devOpFileBytes []byte
    devOpFileBytes, err = this.filesys.GetBytesOfFile(pathToDevOpFile)
    if (nil != err) {
      return
    }

    devOpFile := devOpFile{}
    err = this.yamlCodec.fromYaml(
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
