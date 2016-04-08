package core

import (
  "github.com/dev-op-spec/engine/core/models"
)

func newGetLogForOpRunUseCase(
opRunLogFeed opRunLogFeed,
) getLogForOpRunUseCase {

  return &_getLogForOpRunUseCase{
    opRunLogFeed: opRunLogFeed,
  }

}

type getLogForOpRunUseCase interface {
  Execute(
  opRunId string,
  logChannel chan *models.LogEntry,
  ) (err error)
}

type _getLogForOpRunUseCase struct {
  opRunLogFeed opRunLogFeed
}

func (this _getLogForOpRunUseCase) Execute(
opRunId string,
logChannel chan *models.LogEntry,
) (err error) {

  this.opRunLogFeed.RegisterSubscriber(opRunId, logChannel)

  return
}