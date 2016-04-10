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
  ancestors[]models.OpRunStartedEvent,
  ) (opRunId string, err error)
}

func newRunOpUseCase(
eventStream eventStream,
filesys ports.Filesys,
containerEngine ports.ContainerEngine,
opRunLogFeed opRunLogFeed,
uniqueStringFactory uniqueStringFactory,
yamlCodec yamlCodec,
) runOpUseCase {

  return &_runOpUseCase{
    eventStream:eventStream,
    filesys:filesys,
    containerEngine: containerEngine,
    opRunLogFeed:opRunLogFeed,
    uniqueStringFactory:uniqueStringFactory,
    yamlCodec:yamlCodec,
  }

}

type _runOpUseCase struct {
  eventStream         eventStream
  filesys             ports.Filesys
  containerEngine     ports.ContainerEngine
  opRunLogFeed        opRunLogFeed
  uniqueStringFactory uniqueStringFactory
  yamlCodec           yamlCodec
}

func (this _runOpUseCase) Execute(
req models.RunOpReq,
ancestorOpRunStartedEvents[]models.OpRunStartedEvent,
) (opRunId string, err error) {

  var parentOpRunId string
  if (0 != len(ancestorOpRunStartedEvents)) {
    parentOpRunStartedEvent := ancestorOpRunStartedEvents[len(ancestorOpRunStartedEvents) - 1]
    parentOpRunId = parentOpRunStartedEvent.OpRunId()
  }

  opRunId, err = this.uniqueStringFactory.Construct()
  if (nil != err) {
    return
  }
  opRunStartedEvent := models.NewOpRunStartedEvent(
    time.Now(),
    &parentOpRunId,
    *req.OpUrl,
    opRunId,
  )

  this.eventStream.Publish(opRunStartedEvent)

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

  // guard infinite loop
  for _, ancestorOpRunStartedEvent := range ancestorOpRunStartedEvents {

    if (ancestorOpRunStartedEvent.OpRunOpUrl() == *req.OpUrl) {
      err = errors.New(
        fmt.Sprintf(
          "Unable to run op with url=`%v`. Op recursion is currently not supported.",
          req.OpUrl,
        ),
      )
      return
    }

  }
  ancestorOpRunStartedEvents = append(ancestorOpRunStartedEvents, opRunStartedEvent)

  go func() {

    var opRunExitCode int
    defer func() {

      this.eventStream.Publish(
        models.NewOpRunFinishedEvent(
          time.Now(),
          opRunExitCode,
          opRunId,
        ),
      )

    }()

    if (len(_opFile.SubOps) == 0) {

      logChannel := make(chan *models.LogEntry, 1000)

      // register logChannel as feed publisher
      this.opRunLogFeed.RegisterPublisher(opRunId, logChannel)

      // run op
      opRunExitCode, err = this.containerEngine.RunOp(
        req.OpUrl.Path,
        _opFile.Name,
        logChannel,
      )
      if (nil != err) {
        return
      }

      close(logChannel)

    } else {

      // run child ops
      for _, childOp := range _opFile.SubOps {

        var childOpUrl *models.Url
        childOpUrl, err = models.NewUrl(
          // only support local relative urls for now...
          path.Join(
            filepath.Dir(req.OpUrl.Path),
            childOp.Url),
        )
        if (nil != err) {
          return
        }

        var childOpRunId string
        childOpRunId, err = this.Execute(
          *models.NewRunOpReq(
            childOpUrl,
          ),
          ancestorOpRunStartedEvents,
        )
        if (nil != err) {
          return
        }

        eventChannel := make(chan models.Event, 1000)
        this.eventStream.RegisterSubscriber(eventChannel)

        eventLoop:for {
          var event models.Event
          event = <-eventChannel

          switch event := event.(type) {
          case models.OpRunFinishedEvent:
            // OpRunFinishedEvents

            if (event.OpRunId() == childOpRunId) {
              // our childOpRunId

              if (event.OpRunExitCode() != 0) {
                // if non-zero exit code return immediately

                opRunExitCode = event.OpRunExitCode()
                return

              }else {
                break eventLoop
              }
            }
          default:
          // no op
          }

        }

      }

    }

  }()

  return

}
