package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type addSubOperationUseCase interface {
  Execute(
  req models.AddSubOperationReq,
  ) (err error)
}

func newAddSubOperationUseCase(
filesys ports.Filesys,
pathToOperationFileFactory pathToOperationFileFactory,
yamlCodec yamlCodec,
) addSubOperationUseCase {

  return &_addSubOperationUseCase{
    filesys:filesys,
    pathToOperationFileFactory:pathToOperationFileFactory,
    yamlCodec:yamlCodec,
  }

}

type _addSubOperationUseCase struct {
  filesys                    ports.Filesys
  pathToOperationFileFactory pathToOperationFileFactory
  yamlCodec                  yamlCodec
}

func (this _addSubOperationUseCase) Execute(
req models.AddSubOperationReq,
) (err error) {

  pathToOperationFile := this.pathToOperationFileFactory.Construct(
    req.ProjectUrl,
    &req.OperationName,
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

  newOperationFileSubOperation := operationFileSubOperation{
    Url:req.SubOperationUrl,
  }

  if (len(req.PrecedingSubOperationUrl) > 0) {

    subOperations := []operationFileSubOperation{}

    for _, subOperation := range operationFile.SubOperations {

      subOperations = append(subOperations, subOperation)
      if (subOperation.Url == req.PrecedingSubOperationUrl) {
        subOperations = append(subOperations, newOperationFileSubOperation)
      }

    }

    operationFile.SubOperations = subOperations

  } else {

    operationFile.SubOperations = append(operationFile.SubOperations, newOperationFileSubOperation)

  }

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