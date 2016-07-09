package core

//go:generate counterfeiter -o ./fakeKillOpRunUseCase.go --fake-name fakeKillOpRunUseCase ./ killOpRunUseCase

import (
  "github.com/opctl/engine/core/models"
  "github.com/opctl/engine/core/logging"
  "github.com/opctl/engine/core/ports"
  "time"
)

type killOpRunUseCase interface {
  Execute(
  req models.KillOpRunReq,
  ) (
  correlationId string,
  err error,
  )
}

func newKillOpRunUseCase(
containerEngine ports.ContainerEngine,
eventStream eventStream,
logger logging.Logger,
storage storage,
uniqueStringFactory uniqueStringFactory,
) killOpRunUseCase {

  return _killOpRunUseCase{
    containerEngine:containerEngine,
    eventStream:eventStream,
    logger:logger,
    storage:storage,
    uniqueStringFactory:uniqueStringFactory,
  }

}

type _killOpRunUseCase struct {
  containerEngine     ports.ContainerEngine
  eventStream         eventStream
  logger              logging.Logger
  storage             storage
  uniqueStringFactory uniqueStringFactory
}

func (this _killOpRunUseCase) Execute(
req models.KillOpRunReq,
) (
correlationId string,
err error,
) {

  opRunStartedEvents := this.storage.listOpRunStartedEventsWithRootId(req.OpRunId)
  if (len(opRunStartedEvents) == 0) {
    // guard no runs with provided root
    return
  }

  correlationId = this.uniqueStringFactory.Construct()

  for _, opRunStartedEvent := range opRunStartedEvents {
    go func(opRunStartedEvent models.OpRunStartedEvent) {

      // @TODO: handle failure scenario; currently this can leak docker-compose resources
      err := this.containerEngine.KillOpRun(
        correlationId,
        opRunStartedEvent.OpRunOpUrl(),
        opRunStartedEvent.OpRunId(),
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

      opRunEndedEvent := models.NewOpRunEndedEvent(
        correlationId,
        opRunStartedEvent.OpRunId(),
        models.OpRunOutcomeKilled,
        req.OpRunId,
        time.Now().UTC(),
      )

      this.eventStream.Publish(opRunEndedEvent)

    }(opRunStartedEvent)
  }

  this.storage.deleteOpRunsWithRootId(req.OpRunId)

  return

}
