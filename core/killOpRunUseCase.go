package core

//go:generate counterfeiter -o ./fakeKillOpRunUseCase.go --fake-name fakeKillOpRunUseCase ./ killOpRunUseCase

import (
  "github.com/opspec-io/engine/core/models"
  "github.com/opspec-io/engine/core/logging"
  "github.com/opspec-io/engine/core/ports"
  "time"
  "sort"
  "sync"
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

  // delete here (right away) so other kills or end events don't preempt us
  this.storage.deleteOpRunsWithRootId(req.OpRunId)

  correlationId = this.uniqueStringFactory.Construct()

  // perform async
  go func() {

    // want to operate in reverse of order started
    sort.Sort(OpRunStartedEventDescSorter(opRunStartedEvents))

    var containerKillWaitGroup sync.WaitGroup
    for _, opRunStartedEvent := range opRunStartedEvents {

      containerKillWaitGroup.Add(1)

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

        defer containerKillWaitGroup.Done()

      }(opRunStartedEvent)

    }
    containerKillWaitGroup.Wait()

    // @TODO: make kill events publish in realtime rather than waiting for all resources within the root op run to be reclaimed;
    for _, opRunStartedEvent := range opRunStartedEvents {
      opRunEndedEvent := models.NewOpRunEndedEvent(
        correlationId,
        opRunStartedEvent.OpRunId(),
        models.OpRunOutcomeKilled,
        req.OpRunId,
        time.Now().UTC(),
      )

      this.eventStream.Publish(opRunEndedEvent)
    }
  }()

  return

}
