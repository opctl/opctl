package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type listOperationsUseCase interface {
  Execute(
  projectUrl *models.ProjectUrl,
  ) (operations []models.OperationView, err error)
}

func newListOperationsUseCase(
filesys ports.Filesys,
pathToOperationFileFactory pathToOperationFileFactory,
pathToOperationsDirFactory pathToOperationsDirFactory,
yamlCodec yamlCodec,
) listOperationsUseCase {

  return &_listOperationsUseCase{
    filesys:filesys,
    pathToOperationFileFactory:pathToOperationFileFactory,
    pathToOperationsDirFactory:pathToOperationsDirFactory,
    yamlCodec:yamlCodec,
  }

}

type _listOperationsUseCase struct {
  filesys                    ports.Filesys
  pathToOperationFileFactory pathToOperationFileFactory
  pathToOperationsDirFactory pathToOperationsDirFactory
  yamlCodec                  yamlCodec
}

func (this _listOperationsUseCase) Execute(
projectUrl *models.ProjectUrl,
) (operations []models.OperationView, err error) {

  pathToOperationsDir := this.pathToOperationsDirFactory.Construct(
    projectUrl,
  )

  operationDirNames, err := this.filesys.ListNamesOfChildDirs(
    pathToOperationsDir,
  )
  if (nil != err) {
    return
  }

  for _, operationDirName := range operationDirNames {

    pathToOperationFile := this.pathToOperationFileFactory.Construct(
      projectUrl,
      operationDirName,
    )

    var operationFileBytes []byte
    operationFileBytes, err = this.filesys.GetBytesOfFile(pathToOperationFile)
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

    operationRefViews := []models.OperationRefView{}

    for _, operationSubOperation := range operationFile.SubOperations {

      operationRefView := models.NewOperationRefView(
        operationSubOperation.Name,
      )

      operationRefViews = append(operationRefViews, *operationRefView)

    }

    operationView := models.NewOperationView(
      operationFile.Description,
      operationDirName,
      operationRefViews,
    )

    operations = append(operations, *operationView)

  }

  return

}
