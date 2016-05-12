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

  return &logEntryEmittedEvent{
    correlationId:correlationId,
    logEntryMsg:logEntryMsg,
    logEntryOutputStream:logEntryOutputStream,
    timestamp:timestamp,
  }

}

type LogEntryEmittedEvent interface {
  CorrelationId() string
  LogEntryMsg() string
  LogEntryOutputStream() string
  Timestamp() time.Time
}

type logEntryEmittedEvent struct {
  correlationId        string
  logEntryMsg          string
  logEntryOutputStream string
  timestamp            time.Time
}

func (this logEntryEmittedEvent) MarshalJSON() ([]byte, error) {

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

func (this logEntryEmittedEvent) CorrelationId() string {
  return this.correlationId
}

func (this logEntryEmittedEvent) LogEntryMsg() string {
  return this.logEntryMsg
}

func (this logEntryEmittedEvent) LogEntryOutputStream() string {
  return this.logEntryOutputStream
}

func (this logEntryEmittedEvent) Timestamp() time.Time {
  return this.timestamp
}
