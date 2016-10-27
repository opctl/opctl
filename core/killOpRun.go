package core

import (
  "github.com/opspec-io/sdk-golang/models"
  "time"
  "sort"
  "sync"
)

func (this _api) KillOpRun(
req models.KillOpRunReq,
) (
err error,
) {

  events := this.storage.listOpRunStartedEventsWithRootId(req.OpRunId)
  if (len(events) == 0) {
    // guard no runs with provided root
    return
  }

  // delete here (right away) so other kills or end events don't preempt us
  this.storage.deleteOpRunsWithRootId(req.OpRunId)

  // perform async
  go func() {

    // want to operate in reverse of order started
    sort.Sort(EventDescSorter(events))

    var containerKillWaitGroup sync.WaitGroup
    for _, event := range events {

      containerKillWaitGroup.Add(1)

      go func(event models.Event) {

        this.containerEngine.EnsureContainerRemoved(
          event.OpRunStarted.OpRef,
          event.OpRunStarted.OpRunId,
        )

        defer containerKillWaitGroup.Done()

      }(event)

    }
    containerKillWaitGroup.Wait()

    // @TODO: make kill events publish in realtime rather than waiting for all resources within the root op run to be reclaimed;
    for _, event := range events {
      this.eventStream.Publish(
        models.Event{
          Timestamp:time.Now().UTC(),
          OpRunEnded:&models.OpRunEndedEvent{
            OpRunId:event.OpRunStarted.OpRunId,
            Outcome:models.OpRunOutcomeKilled,
            RootOpRunId:event.OpRunStarted.RootOpRunId,
          },
        },
      )
    }
  }()

  return

}
