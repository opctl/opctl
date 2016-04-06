package logengine

import "github.com/dev-op-spec/engine/core/models"

func NewAddLogEntryForDoRunReq(
logEntry models.LogEntry,
doRunId string,
) *AddLogEntryForDoRunReq {

  return &AddLogEntryForDoRunReq{
    LogEntry:logEntry,
    DoRunId :doRunId,
  }

}

type AddLogEntryForDoRunReq struct {
  LogEntry models.LogEntry
  DoRunId  string
}