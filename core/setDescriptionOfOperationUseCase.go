package core

import (
  "github.com/dev-op-spec/engine/core/ports"
  "github.com/dev-op-spec/engine/core/models"
)

type setDescriptionOfOperationUseCase interface {
  Execute(
  req models.SetDescriptionOfOperationReq,
  ) (err error)
}

func newSetDescriptionOfOperationUseCase(
filesys ports.Filesys,
pathToOperationFileFactory pathToOperationFileFactory,
yamlCodec yamlCodec,
) setDescriptionOfOperationUseCase {

  return &_setDescriptionOfOperationUseCase{
    filesys:filesys,
    pathToOperationFileFactory:pathToOperationFileFactory,
    yamlCodec:yamlCodec,
  }

}

type _setDescriptionOfOperationUseCase struct {
  filesys                  ports.Filesys
  pathToOperationFileFactory pathToOperationFileFactory
  yamlCodec                yamlCodec
}

func (this _setDescriptionOfOperationUseCase) Execute(
req models.SetDescriptionOfOperationReq,
) (err error) {

  pathToOperationFile := this.pathToOperationFileFactory.Construct(
    req.ProjectUrl,
    req.OperationName,
  )

  operationFileBytes, err := this.filesys.GetBytesOfFile(pathToOperationFile)
  if (nil != err) {
    return
  }

  operationFile := operationFile{}
  err = this.yamlCodec.fromYaml(
    operationFileBytes,
    &operationFile,
  )
  if (nil != err) {
    return
  }

  operationFile.Description = req.Description

  operationFileBytes, err = this.yamlCodec.toYaml(&operationFile)
  if (nil != err) {
    return
  }

  err = this.filesys.SaveFile(
    pathToOperationFile,
    operationFileBytes,
  )

  return

}
