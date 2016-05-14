package core

//go:generate counterfeiter -o ./fakeOpRunner.go --fake-name fakeOpRunner ./ opRunner

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/ports"
  "github.com/dev-op-spec/engine/core/logging"
  "path/filepath"
  "fmt"
  "path"
  "time"
  "errors"
  "sync"
)

type opRunner interface {
  Run(
  correlationId string,
  opUrl *models.Url,
  ancestorOpRunStartedEvents[]models.OpRunStartedEvent,
  ) (
  opRunId string,
  err error,
  )

  Kill(
  correlationId string,
  opRunId string,
  ) (
  err error,
  )
}

func newOpRunner(
containerEngine ports.ContainerEngine,
eventStream eventStream,
filesys ports.Filesys,
logger logging.Logger,
uniqueStringFactory uniqueStringFactory,
yamlCodec yamlCodec,
) opRunner {

  return &_opRunner{
    containerEngine: containerEngine,
    eventStream:eventStream,
    filesys:filesys,
    logger:logger,
    uniqueStringFactory:uniqueStringFactory,
    yamlCodec:yamlCodec,
    unfinishedOpRunsByIdMap:make(map[string]models.OpRunStartedEvent),
  }

}

type _opRunner struct {
  containerEngine                ports.ContainerEngine
  eventStream                    eventStream
  filesys                        ports.Filesys
  logger                         logging.Logger
  uniqueStringFactory            uniqueStringFactory
  yamlCodec                      yamlCodec

  unfinishedOpRunsByIdMapRWMutex sync.RWMutex
  unfinishedOpRunsByIdMap        map[string]models.OpRunStartedEvent
}

func (this _opRunner) Run(
correlationId string,
opUrl *models.Url,
ancestorOpRunStartedEvents[]models.OpRunStartedEvent,
) (
opRunId string,
err error,
) {

  var parentOpRunId string
  if (0 != len(ancestorOpRunStartedEvents)) {
    parentOpRunStartedEvent := ancestorOpRunStartedEvents[len(ancestorOpRunStartedEvents) - 1]
    parentOpRunId = parentOpRunStartedEvent.OpRunId()
  }

  opFileBytes, err := this.filesys.GetBytesOfFile(
    filepath.Join(opUrl.Path, "op.yml"),
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

  opRunId, err = this.uniqueStringFactory.Construct()
  if (nil != err) {
    return
  }
  opRunStartedEvent := models.NewOpRunStartedEvent(
    correlationId,
    time.Now(),
    parentOpRunId,
    *opUrl,
    opRunId,
  )

  // guard infinite loop
  for _, ancestorOpRunStartedEvent := range ancestorOpRunStartedEvents {

    if (ancestorOpRunStartedEvent.OpRunOpUrl() == *opUrl) {
      err = errors.New(
        fmt.Sprintf(
          "Unable to run op with url=`%v`. Op recursion is currently not supported.",
          opUrl,
        ),
      )
      return
    }

  }
  ancestorOpRunStartedEvents = append(ancestorOpRunStartedEvents, opRunStartedEvent)

  go func() {

    this.eventStream.Publish(opRunStartedEvent)

    this.unfinishedOpRunsByIdMapRWMutex.Lock()
    this.unfinishedOpRunsByIdMap[opRunId] = opRunStartedEvent
    this.unfinishedOpRunsByIdMapRWMutex.Unlock()

    var opRunExitCode int
    defer func() {

      this.eventStream.Publish(
        models.NewOpRunFinishedEvent(
          correlationId,
          time.Now(),
          opRunExitCode,
          opRunId,
        ),
      )

    }()

    if (len(_opFile.SubOps) == 0) {

      // run op
      opRunExitCode, err = this.containerEngine.RunOp(
        correlationId,
        opUrl.Path,
        _opFile.Name,
        this.logger,
      )
      if (nil != err) {
        return
      }

    } else {

      // run child ops
      for _, childOp := range _opFile.SubOps {

        var childOpUrl *models.Url
        childOpUrl, err = models.NewUrl(
          // only support local relative urls for now...
          path.Join(
            filepath.Dir(opUrl.Path),
            childOp.Url),
        )
        if (nil != err) {
          return
        }

        var childOpRunId string
        childOpRunId, err = this.Run(
          correlationId,
          childOpUrl,
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
            if (event.OpRunId() == childOpRunId) {

              if (event.OpRunExitCode() != 0) {
                // if non-zero exit code return immediately

                opRunExitCode = event.OpRunExitCode()
                return

              } else {
                break eventLoop
              }
            }
          case models.OpRunKilledEvent:
            if (event.OpRunId() == childOpRunId) {
              return
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

func (this _opRunner) Kill(
correlationId string,
opRunId string,
) (
err error,
) {

  this.unfinishedOpRunsByIdMapRWMutex.RLock()
  opRunStartedEvent, isUnfinishedOpRun := this.unfinishedOpRunsByIdMap[opRunId]
  this.unfinishedOpRunsByIdMapRWMutex.RUnlock()

  opRunKilledEvent :=
  models.NewOpRunKilledEvent(
    correlationId,
    time.Now(),
    opRunId,
  )

  this.eventStream.Publish(
    opRunKilledEvent,
  )

  if (isUnfinishedOpRun) {

    this.unfinishedOpRunsByIdMapRWMutex.RLock()
    delete(this.unfinishedOpRunsByIdMap, opRunId)
    this.unfinishedOpRunsByIdMapRWMutex.RUnlock()

    this.containerEngine.KillOpRun(
      correlationId,
      opRunStartedEvent.OpRunOpUrl().Path,
      this.logger,
    )

  }

  return
}
