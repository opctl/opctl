package core

//go:generate counterfeiter -o ./fakeOpRunner.go --fake-name fakeOpRunner ./ opRunner

import (
  "github.com/opspec-io/engine/core/models"
  "github.com/opspec-io/engine/core/ports"
  "github.com/opspec-io/engine/core/logging"
  "github.com/opspec-io/sdk-golang"
  opspecModels "github.com/opspec-io/sdk-golang/models"
  "time"
  "sync"
  "errors"
  "path"
  "path/filepath"
)

type opRunner interface {
  Run(
  correlationId string,
  opArgs map[string]string,
  opBundleUrl string,
  opRunId string,
  parentOpRunId string,
  rootOpRunId string,
  ) (
  err error,
  )
}

func newOpRunner(
containerEngine ports.ContainerEngine,
eventStream eventStream,
logger logging.Logger,
opspecSdk opspec.Sdk,
storage storage,
uniqueStringFactory uniqueStringFactory,
) opRunner {

  return &_opRunner{
    containerEngine: containerEngine,
    eventStream:eventStream,
    logger:logger,
    opspecSdk:opspecSdk,
    storage:storage,
    uniqueStringFactory:uniqueStringFactory,
  }

}

type _opRunner struct {
  containerEngine     ports.ContainerEngine
  eventStream         eventStream
  logger              logging.Logger
  opspecSdk           opspec.Sdk
  storage             storage
  uniqueStringFactory uniqueStringFactory
}

func (this _opRunner) Run(
correlationId string,
opArgs map[string]string,
opBundleUrl string,
opRunId string,
parentOpRunId string,
rootOpRunId string,
) (
err error,
) {

  if (rootOpRunId != opRunId && this.storage.isRootOpRunKilled(rootOpRunId)) {
    // guard: root op run not killed out of band
    return
  }

  op, err := this.opspecSdk.GetOp(
    opBundleUrl,
  )
  if (nil != err) {
    return
  }

  opRunStartedEvent := models.NewOpRunStartedEvent(
    correlationId,
    opBundleUrl,
    opRunId,
    parentOpRunId,
    rootOpRunId,
    time.Now().UTC(),
  )

  this.storage.addOpRunStartedEvent(opRunStartedEvent)

  this.eventStream.Publish(opRunStartedEvent)

  if (nil != op.Run && nil != op.Run.Parallel) {
    err = this.runParallel(
      correlationId,
      opArgs,
      opBundleUrl,
      parentOpRunId,
      rootOpRunId,
      op.Run.Parallel,
    )
  } else if (nil != op.Run && nil != op.Run.Serial) {
    err = this.runSerial(
      correlationId,
      opArgs,
      opBundleUrl,
      parentOpRunId,
      rootOpRunId,
      op.Run.Serial,
    )
  } else {
    err = this.containerEngine.RunOp(
      correlationId,
      opArgs,
      opBundleUrl,
      op.Name,
      opRunId,
      this.logger,
    )
  }

  defer func() {

    if (this.storage.isRootOpRunKilled(rootOpRunId)) {
      // ignore killed op runs; handled by killOpRunUseCase
      return
    }

    var opRunOutcome string
    if (nil != err) {
      opRunOutcome = models.OpRunOutcomeFailed
    } else {
      opRunOutcome = models.OpRunOutcomeSucceeded
    }

    this.eventStream.Publish(
      models.NewOpRunEndedEvent(
        correlationId,
        opRunId,
        opRunOutcome,
        rootOpRunId,
        time.Now().UTC(),
      ),
    )

  }()

  return

}

func (this _opRunner) runSerial(
correlationId string,
opArgs map[string]string,
opBundleUrl string,
parentOpRunId string,
rootOpRunId string,
serialRunDeclaration *opspecModels.SerialRunDeclaration,
) (
err error,
) {

  for _, childRunDeclaration := range *serialRunDeclaration {

    if (nil != childRunDeclaration.Parallel) {
      err = this.runParallel(
        correlationId,
        opArgs,
        opBundleUrl,
        parentOpRunId,
        rootOpRunId,
        childRunDeclaration.Parallel,
      )
    } else if (nil != childRunDeclaration.Serial) {
      err = this.runSerial(
        correlationId,
        opArgs,
        opBundleUrl,
        parentOpRunId,
        rootOpRunId,
        childRunDeclaration.Serial,
      )
    } else {
      err = this.Run(
        correlationId,
        opArgs,
        path.Join(filepath.Dir(opBundleUrl), string(childRunDeclaration.Op)),
        this.uniqueStringFactory.Construct(),
        parentOpRunId,
        rootOpRunId,
      )
    }

    if (nil != err) {
      this.logger(
        models.NewLogEntryEmittedEvent(
          correlationId,
          time.Now().UTC(),
          err.Error(),
          logging.StdErrStream,
        ),
      )
      // end run immediately on any error
      return
    }

  }

  return

}

func (this _opRunner) runParallel(
correlationId string,
opArgs map[string]string,
opBundleUrl string,
parentOpRunId string,
rootOpRunId string,
parallelRunDeclaration *opspecModels.ParallelRunDeclaration,
) (
err error,
) {

  var wg sync.WaitGroup
  var isSubOpRunErrors bool

  // run sub ops
  for _, childRunDeclaration := range *parallelRunDeclaration {
    wg.Add(1)

    var childRunDeclarationError error

    go func(childRunDeclaration opspecModels.RunDeclaration) {
      if (nil != childRunDeclaration.Parallel) {
        childRunDeclarationError = this.runParallel(
          correlationId,
          opArgs,
          opBundleUrl,
          parentOpRunId,
          rootOpRunId,
          childRunDeclaration.Parallel,
        )
      } else if (nil != childRunDeclaration.Serial) {
        childRunDeclarationError = this.runSerial(
          correlationId,
          opArgs,
          opBundleUrl,
          parentOpRunId,
          rootOpRunId,
          childRunDeclaration.Serial,
        )
      } else {
        childRunDeclarationError = this.Run(
          correlationId,
          opArgs,
          path.Join(filepath.Dir(opBundleUrl), string(childRunDeclaration.Op)),
          this.uniqueStringFactory.Construct(),
          parentOpRunId,
          rootOpRunId,
        )
      }

      if (nil != childRunDeclarationError) {
        isSubOpRunErrors = true
        this.logger(
          models.NewLogEntryEmittedEvent(
            correlationId,
            time.Now().UTC(),
            childRunDeclarationError.Error(),
            logging.StdErrStream,
          ),
        )
      }

      defer wg.Done()
    }(childRunDeclaration)
  }

  wg.Wait()

  defer func() {
    if (isSubOpRunErrors) {
      err = errors.New("one or more sub op runs had errors")
    }
  }()

  return

}
