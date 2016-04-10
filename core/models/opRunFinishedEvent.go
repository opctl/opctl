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
    timestamp:timestamp,
    opRunExitCode:opRunExitCode,
    opRunId:opRunId,
  }

}

type OpRunFinishedEvent interface {
  Timestamp() time.Time
  OpRunExitCode() int
  OpRunId() string
  Type() string
}

type opRunFinishedEvent struct {
  timestamp     time.Time
  opRunExitCode int
  opRunId       string
}

func (this opRunFinishedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    Timestamp     time.Time `json:"timestamp"`
    OpRunExitCode int `json:"opRunExitCode"`
    OpRunId       string `json:"opRunId"`
    Type          string `json:"type"`
  }{
    Timestamp:this.Timestamp(),
    OpRunExitCode:this.OpRunExitCode(),
    OpRunId:this.OpRunId(),
    Type:this.Type(),
  }

  return json.Marshal(data)
}

func (this opRunFinishedEvent) Timestamp() time.Time {
  return this.timestamp
}

func (this opRunFinishedEvent) OpRunExitCode() int {
  return this.opRunExitCode
}

func (this opRunFinishedEvent) OpRunId() string {
  return this.opRunId
}

func (this opRunFinishedEvent) Type() string {

  return "OpRunFinished"

}

