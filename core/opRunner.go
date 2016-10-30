package core

//go:generate counterfeiter -o ./fakeOpRunner.go --fake-name fakeOpRunner ./ opRunner

import (
  "github.com/opspec-io/sdk-golang/pkg/bundle"
  "github.com/opspec-io/sdk-golang/models"
  "time"
  "sync"
  "path"
  "path/filepath"
)

type opRunner interface {
  Run(
  opRunId string,
  opRunArgs map[string]string,
  opRef string,
  rootOpRunId string,
  ) (
  err error,
  )
}

func newOpRunner(
containerEngine ContainerEngine,
eventStream eventStream,
bundle bundle.Bundle,
storage storage,
uniqueStringFactory uniqueStringFactory,
) opRunner {

  return &_opRunner{
    containerEngine: containerEngine,
    eventStream:eventStream,
    bundle:bundle,
    storage:storage,
    uniqueStringFactory:uniqueStringFactory,
  }

}

type _opRunner struct {
  containerEngine     ContainerEngine
  eventStream         eventStream
  bundle              bundle.Bundle
  storage             storage
  uniqueStringFactory uniqueStringFactory
}

// Runs an op
func (this _opRunner) Run(
opRunId string,
opRunArgs map[string]string,
opBundleUrl string,
rootOpRunId string,
) (
err error,
) {

  if (rootOpRunId != opRunId && this.storage.isRootOpRunKilled(rootOpRunId)) {
    // guard: root op run not killed out of band
    return
  }

  op, err := this.bundle.GetOp(
    opBundleUrl,
  )
  if (nil != err) {
    return
  }

  opRunStartedEvent := models.Event{
    Timestamp:time.Now().UTC(),
    OpRunStarted:&models.OpRunStartedEvent{
      OpRef:opBundleUrl,
      OpRunId:opRunId,
      RootOpRunId:rootOpRunId,
    },
  }

  this.storage.addOpRunStartedEvent(opRunStartedEvent)

  this.eventStream.Publish(opRunStartedEvent)

  if (nil != op.Run && nil != op.Run.Parallel) {
    err = this.runParallel(
      opRunId,
      opRunArgs,
      opBundleUrl,
      rootOpRunId,
      op.Run.Parallel,
    )
  } else if (nil != op.Run && nil != op.Run.Serial) {
    err = this.runSerial(
      opRunId,
      opRunArgs,
      opBundleUrl,
      rootOpRunId,
      op.Run.Serial,
    )
  } else if (nil != op.Run && "" != op.Run.Op) {
    err = this.Run(
      this.uniqueStringFactory.Construct(),
      opRunArgs,
      path.Join(filepath.Dir(opBundleUrl), string(op.Run.Op)),
      rootOpRunId,
    )
  } else {
    err = this.containerEngine.StartContainer(
      opRunArgs,
      opBundleUrl,
      op.Name,
      opRunId,
      this.eventStream,
      rootOpRunId,
    )
  }

  defer func(err error) {

    if (this.storage.isRootOpRunKilled(rootOpRunId)) {
      // ignore killed op runs; handled by killOpRunUseCase
      return
    }

    var opRunOutcome string
    if (nil != err) {
      opRunOutcome = models.OpRunOutcomeFailed
      this.eventStream.Publish(
        models.Event{
          Timestamp:time.Now().UTC(),
          OpRunEncounteredError: &models.OpRunEncounteredErrorEvent{
            Msg:err.Error(),
            OpRunId:opRunId,
            OpRef:opBundleUrl,
            RootOpRunId:rootOpRunId,
          },
        },
      )
    } else {
      opRunOutcome = models.OpRunOutcomeSucceeded
    }

    this.eventStream.Publish(
      models.Event{
        Timestamp:time.Now().UTC(),
        OpRunEnded:&models.OpRunEndedEvent{
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
rootOpRunId string,
serialRunDeclaration *models.SerialRunDeclaration,
) (
err error,
) {

  for _, childRunDeclaration := range *serialRunDeclaration {

    if (nil != childRunDeclaration.Parallel) {
      err = this.runParallel(
        opRunId,
        opRunArgs,
        opBundleUrl,
        rootOpRunId,
        childRunDeclaration.Parallel,
      )
    } else if (nil != childRunDeclaration.Serial) {
      err = this.runSerial(
        opRunId,
        opRunArgs,
        opBundleUrl,
        rootOpRunId,
        childRunDeclaration.Serial,
      )
    } else {
      err = this.Run(
        this.uniqueStringFactory.Construct(),
        opRunArgs,
        path.Join(filepath.Dir(opBundleUrl), string(childRunDeclaration.Op)),
        rootOpRunId,
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
rootOpRunId string,
parallelRunDeclaration *models.ParallelRunDeclaration,
) (
err error,
) {

  var wg sync.WaitGroup
  var isSubOpRunErrors bool

  // run sub ops
  for _, childRunDeclaration := range *parallelRunDeclaration {
    wg.Add(1)

    var childRunDeclarationError error

    go func(childRunDeclaration models.RunDeclaration) {
      if (nil != childRunDeclaration.Parallel) {
        childRunDeclarationError = this.runParallel(
          opRunId,
          opRunArgs,
          opBundleUrl,
          rootOpRunId,
          childRunDeclaration.Parallel,
        )
      } else if (nil != childRunDeclaration.Serial) {
        childRunDeclarationError = this.runSerial(
          opRunId,
          opRunArgs,
          opBundleUrl,
          rootOpRunId,
          childRunDeclaration.Serial,
        )
      } else {
        childRunDeclarationError = this.Run(
          this.uniqueStringFactory.Construct(),
          opRunArgs,
          path.Join(filepath.Dir(opBundleUrl), string(childRunDeclaration.Op)),
          rootOpRunId,
        )
      }

      if (nil != childRunDeclarationError) {
        isSubOpRunErrors = true
      }

      defer wg.Done()
    }(childRunDeclaration)
  }

  wg.Wait()

  return

}
