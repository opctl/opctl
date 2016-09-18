package core

//go:generate counterfeiter -o ./fakeKillOpRunUseCase.go --fake-name fakeKillOpRunUseCase ./ killOpRunUseCase

import (
  "github.com/opspec-io/sdk-golang/models"
  "github.com/opspec-io/engine/core/ports"
  "time"
  "sort"
  "sync"
)

type killOpRunUseCase interface {
  Execute(
  req models.KillOpRunReq,
  ) (
  err error,
  )
}

func newKillOpRunUseCase(
containerEngine ports.ContainerEngine,
eventStream eventStream,
eventPublisher EventPublisher,
storage storage,
uniqueStringFactory uniqueStringFactory,
) killOpRunUseCase {

  return _killOpRunUseCase{
    containerEngine:containerEngine,
    eventStream:eventStream,
    eventPublisher:eventPublisher,
    storage:storage,
    uniqueStringFactory:uniqueStringFactory,
  }

}

type _killOpRunUseCase struct {
  containerEngine     ports.ContainerEngine
  eventStream         eventStream
  eventPublisher      EventPublisher
  storage             storage
  uniqueStringFactory uniqueStringFactory
}

func (this _killOpRunUseCase) Execute(
req models.KillOpRunReq,
) (
err error,
) {

  opRunStartedEvents := this.storage.listOpRunStartedEventsWithRootId(req.OpRunId)
  if (len(opRunStartedEvents) == 0) {
    // guard no runs with provided root
    return
  }

  // delete here (right away) so other kills or end events don't preempt us
  this.storage.deleteOpRunsWithRootId(req.OpRunId)

  // perform async
  go func() {

    // want to operate in reverse of order started
    sort.Sort(EventDescSorter(opRunStartedEvents))

    var containerKillWaitGroup sync.WaitGroup
    for _, opRunStartedEvent := range opRunStartedEvents {

      containerKillWaitGroup.Add(1)

      go func(opRunStartedEvent models.OpRunStartedEvent) {

        this.containerEngine.EnsureContainerRemoved(
          opRunStartedEvent.OpRef,
          opRunStartedEvent.OpRunId,
        )

        defer containerKillWaitGroup.Done()

      }(opRunStartedEvent)

    }
    containerKillWaitGroup.Wait()

    // @TODO: make kill events publish in realtime rather than waiting for all resources within the root op run to be reclaimed;
    for _, opRunStartedEvent := range opRunStartedEvents {
      this.eventStream.Publish(
        models.Event{
          Timestamp:time.Now().UTC(),
          OpRunEnded:&models.OpRunEndedEvent{
            OpRunId:opRunStartedEvent.OpRunId,
            Outcome:models.OpRunOutcomeKilled,
            RootOpRunId:opRunStartedEvent.RootOpRunId,
          },
        },
      )
    }
  }()

  return

}
