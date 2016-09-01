package core

//go:generate counterfeiter -o ./fakeRunOpUseCase.go --fake-name fakeRunOpUseCase ./ runOpUseCase

import (
  "github.com/opspec-io/engine/core/models"
  "github.com/opspec-io/engine/core/logging"
  "time"
)

type runOpUseCase interface {
  Execute(
  req models.RunOpReq,
  ) (
  opRunId string,
  correlationId string,
  err error,
  )
}

func newRunOpUseCase(
opRunner opRunner,
logger logging.Logger,
pathNormalizer pathNormalizer,
uniqueStringFactory uniqueStringFactory,
) runOpUseCase {

  return &_runOpUseCase{
    opRunner:opRunner,
    logger:logger,
    pathNormalizer:pathNormalizer,
    uniqueStringFactory:uniqueStringFactory,
  }

}

type _runOpUseCase struct {
  opRunner            opRunner
  logger              logging.Logger
  pathNormalizer      pathNormalizer
  uniqueStringFactory uniqueStringFactory
}

func (this _runOpUseCase) Execute(
req models.RunOpReq,
) (
opRunId string,
correlationId string,
err error,
) {

  correlationId = this.uniqueStringFactory.Construct()

  opRunId = this.uniqueStringFactory.Construct()

  go func() {
    err = this.opRunner.Run(
      correlationId,
      req.Args,
      this.pathNormalizer.Normalize(req.OpUrl),
      opRunId,
      "",
      opRunId,
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
  }()

  return

}
