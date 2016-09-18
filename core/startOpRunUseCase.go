package core

//go:generate counterfeiter -o ./fakeStartOpRunUseCase.go --fake-name fakeStartOpRunUseCase ./ startOpRunUseCase

import (
  "github.com/opspec-io/sdk-golang/models"
  "time"
)

type startOpRunUseCase interface {
  Execute(
  req models.StartOpRunReq,
  ) (
  opRunId string,
  err error,
  )
}

func newStartOpRunUseCase(
opRunner opRunner,
eventPublisher EventPublisher,
pathNormalizer pathNormalizer,
uniqueStringFactory uniqueStringFactory,
) startOpRunUseCase {

  return &_startOpRunUseCase{
    opRunner:opRunner,
    eventPublisher:eventPublisher,
    pathNormalizer:pathNormalizer,
    uniqueStringFactory:uniqueStringFactory,
  }

}

type _startOpRunUseCase struct {
  opRunner            opRunner
  eventPublisher      EventPublisher
  pathNormalizer      pathNormalizer
  uniqueStringFactory uniqueStringFactory
}

func (this _startOpRunUseCase) Execute(
req models.StartOpRunReq,
) (
opRunId string,
err error,
) {

  opRunId = this.uniqueStringFactory.Construct()
  rootOpRunId := opRunId // this is root

  go func() {
    err = this.opRunner.Run(
      req.Args,
      this.pathNormalizer.Normalize(req.OpUrl),
      opRunId,
      "",
      rootOpRunId,
    )
    if (nil != err) {
      this.eventPublisher(
        models.Event{
          Timestamp: time.Now().UTC(),
          OpRunEncounteredError: &models.OpRunEncounteredErrorEvent{
            Msg: err.Error(),
            OpRunId:opRunId,
            RootOpRunId:rootOpRunId,
          },
        },
      )
    }
  }()

  return

}
