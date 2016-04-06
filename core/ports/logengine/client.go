package ports

import (
  "github.com/dev-op-spec/engine/core/models"
  "github.com/dev-op-spec/engine/core/models/logengine"
)

type Client interface {
  AddLogEntryForDoRun(
  req logengine.AddLogEntryForDoRunReq,
  ) (err error)

  StreamLogEntriesForDoRun(doRunId string) (logChannel chan *models.LogEntry, err error)
}