package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type addDevOpUseCase interface {
  Execute(
  req models.AddDevOpReq,
  ) (err error)
}

func newAddDevOpUseCase(
filesys ports.Filesys,
pathToDevOpDirFactory pathToDevOpDirFactory,
pathToDevOpFileFactory pathToDevOpFileFactory,
ymlCodec yamlCodec,
containerEngine ports.ContainerEngine,
) addDevOpUseCase {

  return &_addDevOpUseCase{
    filesys:filesys,
    pathToDevOpDirFactory:pathToDevOpDirFactory,
    pathToDevOpFileFactory:pathToDevOpFileFactory,
    ymlCodec:ymlCodec,
    containerEngine:containerEngine,
  }

}

type _addDevOpUseCase struct {
  filesys                ports.Filesys
  pathToDevOpDirFactory  pathToDevOpDirFactory
  pathToDevOpFileFactory pathToDevOpFileFactory
  ymlCodec               yamlCodec
  containerEngine        ports.ContainerEngine
}

func (this _addDevOpUseCase) Execute(
req models.AddDevOpReq,
) (err error) {

  pathToDevOpDir := this.pathToDevOpDirFactory.Construct(
    req.ProjectUrl,
    req.Name,
  )

  err = this.filesys.CreateDir(pathToDevOpDir)
  if (nil != err) {
    return
  }

  var devOpFile = devOpFile{
    Description:req.Description,
  }

  devOpFileBytes, err := this.ymlCodec.toYaml(&devOpFile)
  if (nil != err) {
    return
  }

  pathToDevOpFile := this.pathToDevOpFileFactory.Construct(
    req.ProjectUrl,
    req.Name,
  )

  err = this.filesys.SaveFile(
    pathToDevOpFile,
    devOpFileBytes,
  )
  if (nil != err) {
    return
  }

  err = this.containerEngine.InitDevOp(pathToDevOpDir)

  return

}

