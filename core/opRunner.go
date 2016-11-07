package core

import (
  "github.com/opspec-io/sdk-golang/pkg/bundle"
  "github.com/opspec-io/sdk-golang/pkg/model"
  "time"
  "sync"
  "path"
  "path/filepath"
  "github.com/opspec-io/engine/util/eventing"
  "github.com/opspec-io/engine/pkg/containerengine"
  "github.com/opspec-io/engine/util/uniquestring"
  "github.com/pkg/errors"
)

type opRunner interface {
  Run(
  opRunId string,
  opRunArgs map[string]string,
  opRef string,
  parentOpRunId string,
  ) (
  err error,
  )
}

func newOpRunner(
containerEngine containerengine.ContainerEngine,
eventStream eventing.EventStream,
bundle bundle.Bundle,
opRunRepo opRunRepo,
uniqueStringFactory uniquestring.UniqueStringFactory,
) opRunner {

  return &_opRunner{
    containerEngine: containerEngine,
    eventStream:eventStream,
    bundle:bundle,
    opRunRepo:opRunRepo,
    uniqueStringFactory:uniqueStringFactory,
  }

}

type _opRunner struct {
  containerEngine     containerengine.ContainerEngine
  eventStream         eventing.EventStream
  bundle              bundle.Bundle
  opRunRepo           opRunRepo
  uniqueStringFactory uniquestring.UniqueStringFactory
}

// Runs an op
func (this _opRunner) Run(
opRunId string,
opRunArgs map[string]string,
opRef string,
parentOpRunId string,
) (
err error,
) {

  parentOpRun := this.opRunRepo.tryGet(parentOpRunId)
  if ("" != parentOpRunId && nil == parentOpRun) {
    // parent op killed (we got preempted)
    return
  }

  // determine rootId
  var rootOpRunId string
  if ("" == parentOpRunId) {
    // op is root
    rootOpRunId = opRunId
  } else{
    rootOpRunId = parentOpRun.RootId
  }

  op, err := this.bundle.GetOp(
    opRef,
  )
  if (nil != err) {
    return
  }

  // apply input param defaults
  for _, input := range op.Inputs {
    if _, isArgForInput := opRunArgs[input.Name]; !isArgForInput {
      opRunArgs[input.Name] = input.Default
    }
  }

  opRunStartedEvent := model.Event{
    Timestamp:time.Now().UTC(),
    OpRunStarted:&model.OpRunStartedEvent{
      OpRef:opRef,
      OpRunId:opRunId,
      RootOpRunId:rootOpRunId,
    },
  }
  this.eventStream.Publish(opRunStartedEvent)

  opRun := &model.OpRun{
    Id:opRunId,
    OpRef:opRef,
    ParentId:parentOpRunId,
    RootId:rootOpRunId,
  }
  this.opRunRepo.add(opRun)

  if (nil != op.Run && nil != op.Run.Parallel) {
    err = this.runParallel(
      opRunId,
      opRunArgs,
      opRef,
      parentOpRunId,
      op.Run.Parallel,
    )
  } else if (nil != op.Run && nil != op.Run.Serial) {
    err = this.runSerial(
      opRunId,
      opRunArgs,
      opRef,
      parentOpRunId,
      op.Run.Serial,
    )
  } else if (nil != op.Run && "" != op.Run.Op) {
    err = this.Run(
      this.uniqueStringFactory.Construct(),
      opRunArgs,
      path.Join(filepath.Dir(opRef), string(op.Run.Op)),
      opRunId,
    )
  } else {
    err = this.containerEngine.StartContainer(
      opRunArgs,
      opRef,
      op.Name,
      opRunId,
      this.eventStream,
      rootOpRunId,
    )
  }

  defer func(err error) {

    if ("" != parentOpRunId && nil == this.opRunRepo.tryGet(parentOpRunId)) {
      // guard: parent op killed (we got preempted)
      return
    }

    var opRunOutcome string
    if (nil != err) {
      opRunOutcome = model.OpRunOutcomeFailed
      this.eventStream.Publish(
        model.Event{
          Timestamp:time.Now().UTC(),
          OpRunEncounteredError: &model.OpRunEncounteredErrorEvent{
            Msg:err.Error(),
            OpRunId:opRunId,
            OpRef:opRef,
            RootOpRunId:rootOpRunId,
          },
        },
      )
    } else {
      opRunOutcome = model.OpRunOutcomeSucceeded
    }

    this.eventStream.Publish(
      model.Event{
        Timestamp:time.Now().UTC(),
        OpRunEnded:&model.OpRunEndedEvent{
          OpRunId:opRunId,
          Outcome:opRunOutcome,
          RootOpRunId:rootOpRunId,
        },
      },
    )

  }(err)

  return

}

func (this _opRunner) runSerial(
opRunId string,
opRunArgs map[string]string,
opBundleUrl string,
parentOpRunId string,
serialRunDeclaration *model.SerialRunDeclaration,
) (
err error,
) {

  for _, childRunDeclaration := range *serialRunDeclaration {

    if (nil != childRunDeclaration.Parallel) {
      err = this.runParallel(
        opRunId,
        opRunArgs,
        opBundleUrl,
        parentOpRunId,
        childRunDeclaration.Parallel,
      )
    } else if (nil != childRunDeclaration.Serial) {
      err = this.runSerial(
        opRunId,
        opRunArgs,
        opBundleUrl,
        parentOpRunId,
        childRunDeclaration.Serial,
      )
    } else {
      err = this.Run(
        this.uniqueStringFactory.Construct(),
        opRunArgs,
        path.Join(filepath.Dir(opBundleUrl), string(childRunDeclaration.Op)),
        opRunId,
      )
    }

    if (nil != err) {
      // end run immediately on any error
      return
    }

  }

  return

}

func (this _opRunner) runParallel(
opRunId string,
opRunArgs map[string]string,
opBundleUrl string,
parentOpRunId string,
parallelRunDeclaration *model.ParallelRunDeclaration,
) (
err error,
) {

  var wg sync.WaitGroup
  var isSubOpRunErrors bool

  // run sub ops
  for _, childRunDeclaration := range *parallelRunDeclaration {
    wg.Add(1)

    var childRunDeclarationError error

    go func(childRunDeclaration model.RunDeclaration) {
      if (nil != childRunDeclaration.Parallel) {
        childRunDeclarationError = this.runParallel(
          opRunId,
          opRunArgs,
          opBundleUrl,
          parentOpRunId,
          childRunDeclaration.Parallel,
        )
      } else if (nil != childRunDeclaration.Serial) {
        childRunDeclarationError = this.runSerial(
          opRunId,
          opRunArgs,
          opBundleUrl,
          parentOpRunId,
          childRunDeclaration.Serial,
        )
      } else {
        childRunDeclarationError = this.Run(
          this.uniqueStringFactory.Construct(),
          opRunArgs,
          path.Join(filepath.Dir(opBundleUrl), string(childRunDeclaration.Op)),
          opRunId,
        )
      }

      if (nil != childRunDeclarationError) {
        isSubOpRunErrors = true
      }

      defer wg.Done()
    }(childRunDeclaration)
  }
  wg.Wait()

  if (isSubOpRunErrors) {
    err = errors.New("One or more errors encountered in parallel run block")
  }

  return

}
