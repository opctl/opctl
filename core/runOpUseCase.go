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

type runOpUseCase interface {
  Execute(
  req models.RunOpReq,
  namesOfAlreadyRunOps[]*models.Url,
  ) (opRun models.OpRunDetailedView, err error)
}

func newRunOpUseCase(
filesys ports.Filesys,
containerEngine ports.ContainerEngine,
uniqueStringFactory uniqueStringFactory,
yamlCodec yamlCodec,
) runOpUseCase {

  return &_runOpUseCase{
    filesys:filesys,
    containerEngine: containerEngine,
    uniqueStringFactory:uniqueStringFactory,
    yamlCodec:yamlCodec,
  }

}

type _runOpUseCase struct {
  filesys             ports.Filesys
  containerEngine     ports.ContainerEngine
  uniqueStringFactory uniqueStringFactory
  yamlCodec           yamlCodec
}

func (this _runOpUseCase) Execute(
req models.RunOpReq,
urlsOfAlreadyRunOps[]*models.Url,
) (opRun models.OpRunDetailedView, err error) {

  opRun.StartedAtUnixTime = time.Now().Unix()

  opRun.Id, err = this.uniqueStringFactory.Construct()
  if (nil != err) {
    return
  }

  opRun.OpUrl = req.OpUrl

  opFileBytes, err := this.filesys.GetBytesOfFile(
    filepath.Join(req.OpUrl.Path, "op.yml"),
  )
  if (nil != err) {
    return
  }

  _opFile := &opFile{}
  err = this.yamlCodec.fromYaml(
    opFileBytes,
    _opFile,
  )
  if (nil != err) {
    return
  }

  defer func() {

    opRun.EndedAtUnixTime = time.Now().Unix()

  }()

  // guard infinite loop
  for _, previouslyRunOp := range urlsOfAlreadyRunOps {

    if (*previouslyRunOp == *req.OpUrl) {
      err = errors.New(
        fmt.Sprintf(
          "Unable to run op with url=`%v`. Op recursion is currently not supported.",
          *req.OpUrl,
        ),
      )
      return
    }

  }
  urlsOfAlreadyRunOps = append(urlsOfAlreadyRunOps, req.OpUrl)

  if (len(_opFile.SubOps) == 0) {

    logChannel := make(chan *models.LogEntry, 1000)
    go func() {
      for {
        logEntry := <-logChannel
        fmt.Printf(
          "Timestamp: `%v` | Stream: `%v` | Message: `%v` \n",
          logEntry.Timestamp,
          logEntry.Stream,
          logEntry.Message,
        )
      }
    }()

    // run op
    opRun.ExitCode, err = this.containerEngine.RunOp(
      req.OpUrl.Path,
      *_opFile.Name,
      logChannel,
    )

    if (opRun.ExitCode != 0 || nil != err) {
      return
    }

  } else {

    // run sub ops
    for _, subOp := range _opFile.SubOps {

      var subOpUrl *models.Url
      subOpUrl, err = models.NewUrl(
        // only support local relative urls for now...
        path.Join(
          filepath.Dir(req.OpUrl.Path),
          subOp.Url),
      )
      if (nil != err) {
        return
      }

      var subOpRun models.OpRunDetailedView
      subOpRun, err = this.Execute(
        *models.NewRunOpReq(
          subOpUrl,
        ),
        urlsOfAlreadyRunOps,
      )

      opRun.SubOps = append(
        opRun.SubOps,
        models.NewOpRunSummaryView(
          subOpRun.Id,
          subOpRun.OpUrl,
          subOpRun.StartedAtUnixTime,
          subOpRun.EndedAtUnixTime,
          subOpRun.ExitCode,
        ),
      )

      if (subOpRun.ExitCode != 0 || nil != err) {

        // bubble exit code up
        opRun.ExitCode = subOpRun.ExitCode
        return

      }

    }

  }

  return

}
