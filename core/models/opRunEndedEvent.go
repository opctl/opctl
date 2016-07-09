package models

import (
  "time"
  "encoding/json"
)

func NewOpRunEndedEvent(
correlationId string,
opRunId string,
outcome string,
rootOpRunId string,
timestamp time.Time,
) OpRunEndedEvent {

  return OpRunEndedEvent{
    correlationId:correlationId,
    opRunId:opRunId,
    outcome:outcome,
    rootOpRunId:rootOpRunId,
    timestamp:timestamp,
  }

}

const (
  OpRunOutcomeSucceeded = "SUCCEEDED"
  OpRunOutcomeFailed = "FAILED"
  OpRunOutcomeKilled = "KILLED"
)

type OpRunEndedEvent struct {
  correlationId string
  opRunId       string
  outcome       string
  rootOpRunId   string
  timestamp     time.Time
}

func (this OpRunEndedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    CorrelationId string `json:"correlationId"`
    OpRunId       string `json:"opRunId"`
    Outcome       string `json:"outcome"`
    RootOpRunId   string `json:"rootOpRunId"`
    Timestamp     time.Time `json:"timestamp"`
  }{
    CorrelationId:this.CorrelationId(),
    OpRunId:this.OpRunId(),
    Outcome:this.Outcome(),
    RootOpRunId:this.RootOpRunId(),
    Timestamp:this.Timestamp(),
  }

  return json.Marshal(data)
}

func (this OpRunEndedEvent) CorrelationId() string {
  return this.correlationId
}

func (this OpRunEndedEvent) OpRunId() string {
  return this.opRunId
}

func (this OpRunEndedEvent) Outcome() string {
  return this.outcome
}

func (this OpRunEndedEvent) RootOpRunId() string {
  return this.rootOpRunId
}

func (this OpRunEndedEvent) Timestamp() time.Time {
  return this.timestamp
}
