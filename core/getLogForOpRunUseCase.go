package core

import (
  "github.com/dev-op-spec/engine/core/models"
)

type getLogForOpRunUseCase interface {
  Execute(
  opRunId string,
  ) (logChannel chan *models.LogEntry, err error)
}

type _getLogForOpRunUseCase struct{}

func (this _getLogForOpRunUseCase) Execute(
opRunId string,
) (logChannel chan *models.LogEntry, err error) {
return
}