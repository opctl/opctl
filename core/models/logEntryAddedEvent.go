package models

import (
  "time"
  "encoding/json"
)

func NewLogEntryAddedEvent(
timestamp time.Time,
logEntryMsg string,
logEntryOutputStream string,
) LogEntryAddedEvent {

  return &logEntryAddedEvent{
    timestamp:timestamp,
    logEntryMsg:logEntryMsg,
    logEntryOutputStream:logEntryOutputStream,
  }

}

type LogEntryAddedEvent interface {
  Timestamp() time.Time
  LogEntryMsg() string
  LogEntryOutputStream() string
  Type() string
}

type logEntryAddedEvent struct {
  timestamp            time.Time
  logEntryMsg          string
  logEntryOutputStream string
}

func (this logEntryAddedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    Timestamp            time.Time `json:"timestamp"`
    LogEntryMsg          string `json:"logEntryMsg"`
    LogEntryOutputStream string `json:"logEntryOutputStream"`
    Type                 string `json:"type"`
  }{
    Timestamp:this.Timestamp(),
    LogEntryMsg:this.LogEntryMsg(),
    LogEntryOutputStream:this.LogEntryOutputStream(),
    Type:this.Type(),
  }

  return json.Marshal(data)
}

func (this logEntryAddedEvent) Timestamp() time.Time {
  return this.timestamp
}

func (this logEntryAddedEvent) LogEntryMsg() string {
  return this.logEntryMsg
}

func (this logEntryAddedEvent) LogEntryOutputStream() string {
  return this.logEntryOutputStream
}

func (this logEntryAddedEvent) Type() string {

  return "LogEntryAdded"

}
