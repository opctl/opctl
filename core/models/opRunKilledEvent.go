package models

import (
  "time"
  "encoding/json"
)

func NewOpRunKilledEvent(
correlationId string,
timestamp time.Time,
opRunId string,
) OpRunKilledEvent {

  return &opRunKilledEvent{
    correlationId:correlationId,
    opRunId:opRunId,
    timestamp:timestamp,
  }

}

type OpRunKilledEvent interface {
  CorrelationId() string
  OpRunId() string
  Timestamp() time.Time
}

type opRunKilledEvent struct {
  correlationId string
  opRunId       string
  timestamp     time.Time
}

func (this opRunKilledEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    CorrelationId string `json:"correlationId"`
    OpRunId       string `json:"opRunId"`
    Timestamp     time.Time `json:"timestamp"`
  }{
    CorrelationId:this.CorrelationId(),
    OpRunId:this.OpRunId(),
    Timestamp:this.Timestamp(),
  }

  return json.Marshal(data)
}

func (this opRunKilledEvent) CorrelationId() string {
  return this.correlationId
}

func (this opRunKilledEvent) OpRunId() string {
  return this.opRunId
}

func (this opRunKilledEvent) Timestamp() time.Time {
  return this.timestamp
}
