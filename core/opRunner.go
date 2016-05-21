package core

//go:generate counterfeiter -o ./fakeOpRunner.go --fake-name fakeOpRunner ./ opRunner

import (
  "github.com/open-devops/engine/core/models"
  "github.com/open-devops/engine/core/ports"
  "github.com/open-devops/engine/core/logging"
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
    *opUrl,
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
        correlationId,
        opUrl.Path,
        _opFile.Name,
        this.logger,
      )
      if (nil != err) {
        return
      }

    } else {
      // run sub ops

      for _, subOp := range _opFile.SubOps {

        var subOpUrl *models.Url
        subOpUrl, err = models.NewUrl(
          // only support local relative urls for now...

          path.Join(
            filepath.Dir(opUrl.Path),
            subOp.Url),
        )
        if (nil != err) {
          return
        }

        var subOpRunId string
        subOpRunId, err = this.Run(
          correlationId,
          subOpUrl,
          opRunId,
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
            if (event.OpRunId() == subOpRunId) {

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
opUrl *models.Url,
) (err error) {

  if ("" == parentOpRunId) {
    // handle root op run

    return
  }

  this.unfinishedOpRunsByIdMapRWMutex.RLock()
  parentOpRun := this.unfinishedOpRunsByIdMap[parentOpRunId]
  this.unfinishedOpRunsByIdMapRWMutex.RUnlock()

  if (*opUrl == parentOpRun.OpRunOpUrl()) {
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
        unfinishedOpRun.OpRunOpUrl().Path,
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
