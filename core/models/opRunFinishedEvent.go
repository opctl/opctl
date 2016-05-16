package models

import (
  "time"
  "encoding/json"
)

func NewOpRunFinishedEvent(
correlationId string,
opRunExitCode int,
opRunId string,
rootOpRunId string,
timestamp time.Time,
) OpRunFinishedEvent {

  return OpRunFinishedEvent{
    correlationId:correlationId,
    opRunExitCode:opRunExitCode,
    opRunId:opRunId,
    rootOpRunId:rootOpRunId,
    timestamp:timestamp,
  }

}

type OpRunFinishedEvent struct {
  correlationId string
  opRunExitCode int
  opRunId       string
  rootOpRunId   string
  timestamp     time.Time
}

func (this OpRunFinishedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    CorrelationId string `json:"correlationId"`
    OpRunExitCode int `json:"opRunExitCode"`
    OpRunId       string `json:"opRunId"`
    RootOpRunId   string `json:"rootOpRunId"`
    Timestamp     time.Time `json:"timestamp"`
  }{
    CorrelationId:this.CorrelationId(),
    OpRunExitCode:this.OpRunExitCode(),
    OpRunId:this.OpRunId(),
    RootOpRunId:this.RootOpRunId(),
    Timestamp:this.Timestamp(),
  }

  return json.Marshal(data)
}

func (this OpRunFinishedEvent) CorrelationId() string {
  return this.correlationId
}

func (this OpRunFinishedEvent) OpRunExitCode() int {
  return this.opRunExitCode
}

func (this OpRunFinishedEvent) OpRunId() string {
  return this.opRunId
}

func (this OpRunFinishedEvent) RootOpRunId() string {
  return this.rootOpRunId
}

func (this OpRunFinishedEvent) Timestamp() time.Time {
  return this.timestamp
}
