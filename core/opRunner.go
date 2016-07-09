package core

//go:generate counterfeiter -o ./fakeOpRunner.go --fake-name fakeOpRunner ./ opRunner

import (
  "github.com/opctl/engine/core/models"
  "github.com/opctl/engine/core/ports"
  "github.com/opctl/engine/core/logging"
  "github.com/opspec-io/sdk-golang"
  opspecModels "github.com/opspec-io/sdk-golang/models"
  "time"
  "sync"
  "path"
  "path/filepath"
  "errors"
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

  if (len(op.SubOps) != 0) {
    err = this.runSubOps(
      correlationId,
      opArgs,
      opBundleUrl,
      parentOpRunId,
      rootOpRunId,
      op.SubOps,
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

func (this _opRunner) runSubOps(
correlationId string,
opArgs map[string]string,
opBundleUrl string,
parentOpRunId string,
rootOpRunId string,
subOps []opspecModels.SubOpView,
) (
err error,
) {

  var wg sync.WaitGroup
  var isSubOpRunErrors bool

  // run sub ops
  for _, subOp := range subOps {

    wg.Add(1)

    // currently only support embedded sub ops
    subOpBundleUrl := path.Join(
      filepath.Dir(opBundleUrl),
      subOp.Url,
    )

    subOpRunId := this.uniqueStringFactory.Construct()

    if subOp.IsParallel {
      // handle parallel
      go func() {
        subOpRunErr := this.Run(
          correlationId,
          opArgs,
          subOpBundleUrl,
          subOpRunId,
          parentOpRunId,
          rootOpRunId,
        )
        if (nil != subOpRunErr) {
          isSubOpRunErrors = true
          this.logger(
            models.NewLogEntryEmittedEvent(
              correlationId,
              time.Now().UTC(),
              subOpRunErr.Error(),
              logging.StdErrStream,
            ),
          )
        }

        wg.Done()
      }()
    } else {
      // handle synchronous
      subOpRunErr := this.Run(
        correlationId,
        opArgs,
        subOpBundleUrl,
        subOpRunId,
        parentOpRunId,
        rootOpRunId,
      )
      if (nil != subOpRunErr) {
        isSubOpRunErrors = true
        this.logger(
          models.NewLogEntryEmittedEvent(
            correlationId,
            time.Now().UTC(),
            subOpRunErr.Error(),
            logging.StdErrStream,
          ),
        )
      }

      wg.Done()
    }

  }

  wg.Wait()

  defer func() {
    if (isSubOpRunErrors) {
      err = errors.New("one or more sub op runs had errors")
    }
  }()

  return

}
