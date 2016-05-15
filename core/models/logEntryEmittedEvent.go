package models

import (
  "time"
  "encoding/json"
)

func NewLogEntryEmittedEvent(
correlationId string,
timestamp time.Time,
logEntryMsg string,
logEntryOutputStream string,
) LogEntryEmittedEvent {

  return LogEntryEmittedEvent{
    correlationId:correlationId,
    logEntryMsg:logEntryMsg,
    logEntryOutputStream:logEntryOutputStream,
    timestamp:timestamp,
  }

}

type LogEntryEmittedEvent struct {
  correlationId        string
  logEntryMsg          string
  logEntryOutputStream string
  timestamp            time.Time
}

func (this LogEntryEmittedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    CorrelationId        string `json:"correlationId"`
    LogEntryMsg          string `json:"logEntryMsg"`
    LogEntryOutputStream string `json:"logEntryOutputStream"`
    Timestamp            time.Time `json:"timestamp"`
  }{
    CorrelationId:this.CorrelationId(),
    LogEntryMsg:this.LogEntryMsg(),
    LogEntryOutputStream:this.LogEntryOutputStream(),
    Timestamp:this.Timestamp(),
  }

  return json.Marshal(data)
}

func (this LogEntryEmittedEvent) CorrelationId() string {
  return this.correlationId
}

func (this LogEntryEmittedEvent) LogEntryMsg() string {
  return this.logEntryMsg
}

func (this LogEntryEmittedEvent) LogEntryOutputStream() string {
  return this.logEntryOutputStream
}

func (this LogEntryEmittedEvent) Timestamp() time.Time {
  return this.timestamp
}
