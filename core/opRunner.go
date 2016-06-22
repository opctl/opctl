package core

//go:generate counterfeiter -o ./fakeOpRunner.go --fake-name fakeOpRunner ./ opRunner

import (
  "github.com/opctl/engine/core/models"
  "github.com/opctl/engine/core/ports"
  "github.com/opctl/engine/core/logging"
  "github.com/opspec-io/sdk-golang"
  "path/filepath"
  "fmt"
  "path"
  "time"
  "errors"
  "sync"
)

type opRunner interface {
  Run(
  args map[string]string,
  correlationId string,
  opUrl string,
  parentOpRunId string,
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
logger logging.Logger,
opspecSdk opspec.Sdk,
uniqueStringFactory uniqueStringFactory,
) opRunner {

  return &_opRunner{
    containerEngine: containerEngine,
    eventStream:eventStream,
    logger:logger,
    opspecSdk:opspecSdk,
    uniqueStringFactory:uniqueStringFactory,
    unfinishedOpRunsByIdMap:make(map[string]models.OpRunStartedEvent),
  }

}

type _opRunner struct {
  containerEngine                ports.ContainerEngine
  eventStream                    eventStream
  logger                         logging.Logger
  opspecSdk                      opspec.Sdk
  uniqueStringFactory            uniqueStringFactory

  unfinishedOpRunsByIdMapRWMutex sync.RWMutex
  unfinishedOpRunsByIdMap        map[string]models.OpRunStartedEvent
}

func (this _opRunner) Run(
args map[string]string,
correlationId string,
opUrl string,
parentOpRunId string,
) (
opRunId string,
err error,
) {

  err = this.guardOpRunNotRecursive(
    parentOpRunId,
    opUrl,
  )
  if (nil != err) {
    return
  }

  _opFile, err := this.opspecSdk.GetOp(
    opUrl,
  )
  if (nil != err) {
    return
  }

  opRunId, err = this.uniqueStringFactory.Construct()
  if (nil != err) {
    return
  }

  var rootOpRunId string
  if ("" == parentOpRunId) {
    // handle root op run

    rootOpRunId = opRunId

  } else {
    // handle sub op run

    this.unfinishedOpRunsByIdMapRWMutex.RLock()
    rootOpRunId = this.unfinishedOpRunsByIdMap[parentOpRunId].RootOpRunId()
    this.unfinishedOpRunsByIdMapRWMutex.RUnlock()

  }

  opRunStartedEvent := models.NewOpRunStartedEvent(
    correlationId,
    opUrl,
    opRunId,
    parentOpRunId,
    rootOpRunId,
    time.Now().UTC(),
  )

  this.eventStream.Publish(opRunStartedEvent)

  this.unfinishedOpRunsByIdMapRWMutex.Lock()
  this.unfinishedOpRunsByIdMap[opRunId] = opRunStartedEvent
  this.unfinishedOpRunsByIdMapRWMutex.Unlock()

  go func() {

    var opRunExitCode int
    defer func() {

      this.unfinishedOpRunsByIdMapRWMutex.RLock()
      _, isUnfinishedOpRun := this.unfinishedOpRunsByIdMap[opRunId]
      this.unfinishedOpRunsByIdMapRWMutex.RUnlock()

      if (isUnfinishedOpRun) {

        this.unfinishedOpRunsByIdMapRWMutex.Lock()
        delete(this.unfinishedOpRunsByIdMap, opRunId)
        this.unfinishedOpRunsByIdMapRWMutex.Unlock()

        this.eventStream.Publish(
          models.NewOpRunFinishedEvent(
            correlationId,
            opRunExitCode,
            opRunId,
            rootOpRunId,
            time.Now().UTC(),
          ),
        )

      }

    }()

    if (len(_opFile.SubOps) == 0) {
      // run op

      opRunExitCode, err = this.containerEngine.RunOp(
        args,
        correlationId,
        opUrl,
        _opFile.Name,
        this.logger,
      )
      if (nil != err) {
        this.logger(
          models.NewLogEntryEmittedEvent(
            correlationId,
            time.Now().UTC(),
            err.Error(),
            logging.StdErrStream,
          ),
        )
      }

    } else {
      // run sub ops

      for _, subOp := range _opFile.SubOps {

        subOpUrl := path.Join(
          filepath.Dir(opUrl),
          subOp.Name,
        )

        eventChannel := make(chan models.Event)
        this.eventStream.RegisterSubscriber(eventChannel)
        defer this.eventStream.UnregisterSubscriber(eventChannel)

        var subOpRunId string
        subOpRunId, err = this.Run(
          args,
          correlationId,
          subOpUrl,
          opRunId,
        )
        if (nil != err) {
          this.logger(
            models.NewLogEntryEmittedEvent(
              correlationId,
              time.Now().UTC(),
              err.Error(),
              logging.StdErrStream,
            ),
          )
        }

        eventLoop:for {
          event := <-eventChannel

          switch event := event.(type) {
          case models.OpRunFinishedEvent:
            if (event.OpRunId() == subOpRunId) {

              this.eventStream.UnregisterSubscriber(eventChannel)

              if (event.OpRunExitCode() != 0) {
                // if non-zero exit code return immediately

                opRunExitCode = event.OpRunExitCode()
                return

              } else {
                break eventLoop
              }
            }
          default: // no op
          }

        }

      }

    }

  }()

  return

}

func (this _opRunner) guardOpRunNotRecursive(
parentOpRunId string,
opUrl string,
) (err error) {

  if ("" == parentOpRunId) {
    // handle root op run

    return
  }

  this.unfinishedOpRunsByIdMapRWMutex.RLock()
  parentOpRun := this.unfinishedOpRunsByIdMap[parentOpRunId]
  this.unfinishedOpRunsByIdMapRWMutex.RUnlock()

  if (opUrl == parentOpRun.OpRunOpUrl()) {
    // handle infinite recursion

    return errors.New(
      fmt.Sprintf(
        "Unable to run op with url=`%v`. Found op recursion.",
        opUrl,
      ),
    )
  }

  return this.guardOpRunNotRecursive(
    parentOpRun.ParentOpRunId(),
    opUrl,
  )

}

func (this _opRunner) Kill(
correlationId string,
opRunId string,
) (
err error,
) {

  // get rootOpRunId
  this.unfinishedOpRunsByIdMapRWMutex.RLock()
  rootOpRunId := this.unfinishedOpRunsByIdMap[opRunId].RootOpRunId()
  this.unfinishedOpRunsByIdMapRWMutex.RUnlock()

  this.unfinishedOpRunsByIdMapRWMutex.Lock()
  for _, unfinishedOpRun := range this.unfinishedOpRunsByIdMap {

    delete(this.unfinishedOpRunsByIdMap, unfinishedOpRun.OpRunId())

    this.eventStream.Publish(
      models.NewOpRunKilledEvent(
        correlationId,
        unfinishedOpRun.OpRunId(),
        rootOpRunId,
        time.Now(),
      ),
    )

    go func() {

      // @TODO: handle failure scenario; currently this can leak docker-compose resources
      err := this.containerEngine.KillOpRun(
        correlationId,
        unfinishedOpRun.OpRunOpUrl(),
        this.logger,
      )
      if (nil != err) {
        this.logger(
          models.NewLogEntryEmittedEvent(
            correlationId,
            time.Now().UTC(),
            err.Error(),
            logging.StdErrStream,
          ),
        )
      }

      this.eventStream.Publish(
        models.NewOpRunFinishedEvent(
          correlationId,
          138,
          opRunId,
          rootOpRunId,
          time.Now().UTC(),
        ),
      )

    }()

  }
  this.unfinishedOpRunsByIdMapRWMutex.Unlock()

  return
}
