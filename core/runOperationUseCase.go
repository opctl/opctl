package core

import (
  "errors"
  "time"
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
  "path/filepath"
)

type runOperationUseCase interface {
  Execute(
  req models.RunOperationReq,
  namesOfAlreadyRunOperations[]string,
  ) (operationRun models.OperationRunDetailedView, err error)
}

func newRunOperationUseCase(
filesys ports.Filesys,
pathToOperationDirFactory pathToOperationDirFactory,
pathToOperationFileFactory pathToOperationFileFactory,
containerEngine ports.ContainerEngine,
uniqueStringFactory uniqueStringFactory,
yamlCodec yamlCodec,
) runOperationUseCase {

  return &_runOperationUseCase{
    filesys:filesys,
    pathToOperationDirFactory:pathToOperationDirFactory,
    pathToOperationFileFactory:pathToOperationFileFactory,
    containerEngine: containerEngine,
    uniqueStringFactory:uniqueStringFactory,
    yamlCodec:yamlCodec,
  }

}

type _runOperationUseCase struct {
  filesys                    ports.Filesys
  pathToOperationDirFactory  pathToOperationDirFactory
  pathToOperationFileFactory pathToOperationFileFactory
  containerEngine            ports.ContainerEngine
  uniqueStringFactory        uniqueStringFactory
  yamlCodec                  yamlCodec
}

func (this _runOperationUseCase) Execute(
req models.RunOperationReq,
namesOfAlreadyRunOperations[]string,
) (operationRun models.OperationRunDetailedView, err error) {

  pathToOperationFile := this.pathToOperationFileFactory.Construct(
    req.ProjectUrl,
    req.OperationName,
  )

  operationRun.StartedAtUnixTime = time.Now().Unix()

  operationRun.Id, err = this.uniqueStringFactory.Construct()
  if (nil != err) {
    return
  }

  operationRun.OperationName = req.OperationName

  operationFileBytes, err := this.filesys.GetBytesOfFile(pathToOperationFile)
  if (nil != err) {
    return
  }

  _operationFile := &operationFile{}
  err = this.yamlCodec.fromYaml(
    operationFileBytes,
    _operationFile,
  )
  if (nil != err) {
    return
  }

  defer func() {

    operationRun.EndedAtUnixTime = time.Now().Unix()

  }()

  // guard infinite loop
  for _, previouslyRunOperation := range namesOfAlreadyRunOperations {

    if (previouslyRunOperation == req.OperationName) {
      err = errors.New("Unable to run operation with name=`" + req.OperationName +
      "`. Operations cannot call themselves recursively.")
      return
    }

  }
  namesOfAlreadyRunOperations = append(namesOfAlreadyRunOperations, req.OperationName)

  if (len(_operationFile.SubOperations) == 0) {

    // run operation

    var containerEngineOperationRun models.OperationRunDetailedView
    containerEngineOperationRun, err = this.containerEngine.RunOperation(
      filepath.Dir(pathToOperationFile),
    )

    if (containerEngineOperationRun.ExitCode != 0 || nil != err) {

      // bubble exit code up
      operationRun.ExitCode = containerEngineOperationRun.ExitCode
      return

    }

  } else {

    // run sub operations
    for _, subOperation := range _operationFile.SubOperations {

      pathToSubOperationFile := this.pathToOperationFileFactory.Construct(
        req.ProjectUrl,
        subOperation.Name,
      )

      var subOperationFileBytes []byte
      subOperationFileBytes, err = this.filesys.GetBytesOfFile(pathToSubOperationFile)
      if (nil != err) {
        return
      }

      subOperationFile := operationFile{}
      err = this.yamlCodec.fromYaml(
        subOperationFileBytes,
        &subOperationFile,
      )
      if (nil != err) {
        return
      }

      var subOperationRun models.OperationRunDetailedView
      subOperationRun, err = this.Execute(
        *models.NewRunOperationReq(
          req.ProjectUrl,
          subOperation.Name,
        ),
        namesOfAlreadyRunOperations,
      )

      operationRun.SubOperations = append(
        operationRun.SubOperations,
        models.NewOperationRunSummaryView(
          subOperationRun.Id,
          subOperationRun.OperationName,
          subOperationRun.StartedAtUnixTime,
          subOperationRun.EndedAtUnixTime,
          subOperationRun.ExitCode,
        ),
      )

      if (subOperationRun.ExitCode != 0 || nil != err) {

        // bubble exit code up
        operationRun.ExitCode = subOperationRun.ExitCode
        return

      }

    }

  }

  return

}
