package core

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
)

type listOperationsUseCase interface {
  Execute(
  projectUrl *models.Url,
  ) (operations []models.OperationDetailedView, err error)
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
projectUrl *models.Url,
) (operations []models.OperationDetailedView, err error) {

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
      &operationDirName,
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

    operationSummaryViews := []models.OperationSummaryView{}

    for _, operationSubOperation := range operationFile.SubOperations {

      operationSummaryView := models.NewOperationSummaryView(
        operationSubOperation.Url,
      )

      operationSummaryViews = append(operationSummaryViews, *operationSummaryView)

    }

    operationDetailedView := models.NewOperationDetailedView(
      operationFile.Description,
      operationDirName,
      operationSummaryViews,
    )

    operations = append(operations, *operationDetailedView)

  }

  return

}
