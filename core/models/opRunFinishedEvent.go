package models

import (
  "time"
  "encoding/json"
)

func NewOpRunFinishedEvent(
correlationId string,
timestamp time.Time,
opRunExitCode int,
opRunId string,
) OpRunFinishedEvent {

  return &opRunFinishedEvent{
    correlationId:correlationId,
    opRunExitCode:opRunExitCode,
    opRunId:opRunId,
    timestamp:timestamp,
  }

}

type OpRunFinishedEvent interface {
  CorrelationId() string
  OpRunExitCode() int
  OpRunId() string
  Timestamp() time.Time
}

type opRunFinishedEvent struct {
  correlationId string
  opRunExitCode int
  opRunId       string
  timestamp     time.Time
}

func (this opRunFinishedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    CorrelationId string `json:"correlationId"`
    OpRunExitCode int `json:"opRunExitCode"`
    OpRunId       string `json:"opRunId"`
    Timestamp     time.Time `json:"timestamp"`
  }{
    CorrelationId:this.CorrelationId(),
    OpRunExitCode:this.OpRunExitCode(),
    OpRunId:this.OpRunId(),
    Timestamp:this.Timestamp(),
  }

  return json.Marshal(data)
}

func (this opRunFinishedEvent) CorrelationId() string {
  return this.correlationId
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
