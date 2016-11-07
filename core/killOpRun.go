package core

import (
  "github.com/opspec-io/sdk-golang/pkg/model"
  "time"
  "sync"
)

func (this _core) KillOpRun(
req model.KillOpRunReq,
) {

  this.KillOpRunsRecursively(req.OpRunId)

}

func (this _core) KillOpRunsRecursively(sourceOpRunId string) {

  var waitGroup sync.WaitGroup

  opRun := this.opRunRepo.tryGet(sourceOpRunId)
  // guard opRun found
  if (nil == opRun) {
    return
  }

  // order of the following matters (hence the numbering)
  // 1) delete from storage (done first to ensure we're not preempted by another kill request or run completion)
  this.opRunRepo.delete(sourceOpRunId)

  // 2) recover resources
  go func() {
    waitGroup.Add(1)
    this.containerEngine.EnsureContainerRemoved(
      opRun.OpRef,
      opRun.Id,
    )
    defer waitGroup.Done()
  }()

  // 3) kill children
  for _, childOpRun := range this.opRunRepo.listWithParentId(sourceOpRunId) {
    go func(childOpRunId string) {
      waitGroup.Add(1)
      this.KillOpRunsRecursively(childOpRunId)
      defer waitGroup.Done()
    }(childOpRun.Id)
  }

  // 4) wait for 2 & 3
  waitGroup.Wait()

  // 5) send OpRunEndedEvent last to ensure OpRunEndedEvent's send in reverse order of OpRunStartedEvent's
  this.eventStream.Publish(
    model.Event{
      Timestamp:time.Now().UTC(),
      OpRunEnded:&model.OpRunEndedEvent{
        OpRunId:opRun.Id,
        Outcome:model.OpRunOutcomeKilled,
      },
    },
  )

}
