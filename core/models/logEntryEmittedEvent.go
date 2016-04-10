package models

import (
  "time"
  "encoding/json"
)

func NewLogEntryEmittedEvent(
timestamp time.Time,
logEntryMsg string,
logEntryOutputStream string,
) LogEntryEmittedEvent {

  return &logEntryEmittedEvent{
    logEntryMsg:logEntryMsg,
    logEntryOutputStream:logEntryOutputStream,
    timestamp:timestamp,
  }

}

type LogEntryEmittedEvent interface {
  LogEntryMsg() string
  LogEntryOutputStream() string
  Timestamp() time.Time
  Type() string
}

type logEntryEmittedEvent struct {
  logEntryMsg          string
  logEntryOutputStream string
  timestamp            time.Time
}

func (this logEntryEmittedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    LogEntryMsg          string `json:"logEntryMsg"`
    LogEntryOutputStream string `json:"logEntryOutputStream"`
    Timestamp            time.Time `json:"timestamp"`
    Type                 string `json:"type"`
  }{
    LogEntryMsg:this.LogEntryMsg(),
    LogEntryOutputStream:this.LogEntryOutputStream(),
    Timestamp:this.Timestamp(),
    Type:this.Type(),
  }

  return json.Marshal(data)
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

func (this logEntryEmittedEvent) Type() string {

  return "LogEntryAdded"

}
