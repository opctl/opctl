package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type addOperationUseCase interface {
  Execute(
  req models.AddOperationReq,
  ) (err error)
}

func newAddOperationUseCase(
filesys ports.Filesys,
pathToOperationDirFactory pathToOperationDirFactory,
pathToOperationFileFactory pathToOperationFileFactory,
yamlCodec yamlCodec,
) addOperationUseCase {

  return &_addOperationUseCase{
    filesys:filesys,
    pathToOperationDirFactory:pathToOperationDirFactory,
    pathToOperationFileFactory:pathToOperationFileFactory,
    yamlCodec:yamlCodec,
  }

}

type _addOperationUseCase struct {
  filesys                    ports.Filesys
  pathToOperationDirFactory  pathToOperationDirFactory
  pathToOperationFileFactory pathToOperationFileFactory
  yamlCodec                  yamlCodec
}

func (this _addOperationUseCase) Execute(
req models.AddOperationReq,
) (err error) {

  pathToOperationDir := this.pathToOperationDirFactory.Construct(
    req.ProjectUrl,
    &req.Name,
  )

  err = this.filesys.CreateDir(pathToOperationDir)
  if (nil != err) {
    return
  }

  var operationFile = operationFile{
    Description:req.Description,
    Name:&req.Name,
  }

  operationFileBytes, err := this.yamlCodec.toYaml(&operationFile)
  if (nil != err) {
    return
  }

  pathToOperationFile := this.pathToOperationFileFactory.Construct(
    req.ProjectUrl,
    &req.Name,
  )

  err = this.filesys.SaveFile(
    pathToOperationFile,
    operationFileBytes,
  )

  return

}
