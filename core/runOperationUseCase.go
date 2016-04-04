package core

import (
  "errors"
  "time"
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
  "path/filepath"
  "fmt"
  "path"
)

type runOperationUseCase interface {
  Execute(
  req models.RunOperationReq,
  namesOfAlreadyRunOperations[]*models.Url,
  ) (operationRun models.OperationRunDetailedView, err error)
}

func newRunOperationUseCase(
filesys ports.Filesys,
containerEngine ports.ContainerEngine,
uniqueStringFactory uniqueStringFactory,
yamlCodec yamlCodec,
) runOperationUseCase {

  return &_runOperationUseCase{
    filesys:filesys,
    containerEngine: containerEngine,
    uniqueStringFactory:uniqueStringFactory,
    yamlCodec:yamlCodec,
  }

}

type _runOperationUseCase struct {
  filesys             ports.Filesys
  containerEngine     ports.ContainerEngine
  uniqueStringFactory uniqueStringFactory
  yamlCodec           yamlCodec
}

func (this _runOperationUseCase) Execute(
req models.RunOperationReq,
urlsOfAlreadyRunOperations[]*models.Url,
) (operationRun models.OperationRunDetailedView, err error) {

  operationRun.StartedAtUnixTime = time.Now().Unix()

  operationRun.Id, err = this.uniqueStringFactory.Construct()
  if (nil != err) {
    return
  }

  operationRun.OperationUrl = req.OperationUrl

  operationFileBytes, err := this.filesys.GetBytesOfFile(
    filepath.Join(req.OperationUrl.Path, "operation.yml"),
  )
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
  for _, previouslyRunOperation := range urlsOfAlreadyRunOperations {

    if (*previouslyRunOperation == *req.OperationUrl) {
      err = errors.New(
        fmt.Sprintf(
          "Unable to run operation with url=`%v`. Operation recursion is currently not supported.",
          *req.OperationUrl,
        ),
      )
      return
    }

  }
  urlsOfAlreadyRunOperations = append(urlsOfAlreadyRunOperations, req.OperationUrl)

  if (len(_operationFile.SubOperations) == 0) {

    // run operation
    operationRun.ExitCode, err = this.containerEngine.RunOperation(
      req.OperationUrl.Path,
    )

    if (operationRun.ExitCode != 0 || nil != err) {
      return
    }

  } else {

    // run sub operations
    for _, subOperation := range _operationFile.SubOperations {

      var subOperationUrl *models.Url
      subOperationUrl, err = models.NewUrl(
        // only support local relative urls for now...
        path.Join(
          filepath.Dir(req.OperationUrl.Path),
          subOperation.Url),
      )
      if (nil != err) {
        return
      }

      var subOperationRun models.OperationRunDetailedView
      subOperationRun, err = this.Execute(
        *models.NewRunOperationReq(
          subOperationUrl,
        ),
        urlsOfAlreadyRunOperations,
      )

      operationRun.SubOperations = append(
        operationRun.SubOperations,
        models.NewOperationRunSummaryView(
          subOperationRun.Id,
          subOperationRun.OperationUrl,
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
