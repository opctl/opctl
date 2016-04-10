package models

import (
  "time"
  "encoding/json"
)

func NewOpRunFinishedEvent(
timestamp time.Time,
opRunExitCode int,
opRunId string,
) OpRunFinishedEvent {

  return &opRunFinishedEvent{
    opRunExitCode:opRunExitCode,
    opRunId:opRunId,
    timestamp:timestamp,
  }

}

type OpRunFinishedEvent interface {
  OpRunExitCode() int
  OpRunId() string
  Timestamp() time.Time
  Type() string
}

type opRunFinishedEvent struct {
  opRunExitCode int
  opRunId       string
  timestamp     time.Time
}

func (this opRunFinishedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    OpRunExitCode int `json:"opRunExitCode"`
    OpRunId       string `json:"opRunId"`
    Timestamp     time.Time `json:"timestamp"`
    Type          string `json:"type"`
  }{
    OpRunExitCode:this.OpRunExitCode(),
    OpRunId:this.OpRunId(),
    Timestamp:this.Timestamp(),
    Type:this.Type(),
  }

  return json.Marshal(data)
}

func (this opRunFinishedEvent) OpRunExitCode() int {
  return this.opRunExitCode
}

func (this opRunFinishedEvent) OpRunId() string {
  return this.opRunId
}

func (this opRunFinishedEvent) Timestamp() time.Time {
  return this.timestamp
}

func (this opRunFinishedEvent) Type() string {

  return "OpRunFinished"

}

