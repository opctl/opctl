package dockercompose

//go:generate counterfeiter -o ./fakeKillOpRunUseCase.go --fake-name fakeKillOpRunUseCase ./ killOpRunUseCase

import (
  "github.com/dev-op-spec/engine/core/logging"
)

type killOpRunUseCase interface {
  Execute(
  correlationId string,
  pathToOpDir string,
  logger logging.Logger,
  ) (err error)
}

func newKillOpRunUseCase(
opRunResourceFlusher opRunResourceFlusher,
) killOpRunUseCase {

  return &_killOpRunUseCase{
    opRunResourceFlusher: opRunResourceFlusher,
  }

}

type _killOpRunUseCase struct {
  opRunResourceFlusher opRunResourceFlusher
}

func (this _killOpRunUseCase) Execute(
correlationId string,
pathToOpDir string,
logger logging.Logger,
) (err error) {

  this.opRunResourceFlusher.flush(
    correlationId,
    pathToOpDir,
    logger,
  )

  return

}
