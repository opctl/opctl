package models

import (
  "time"
  "encoding/json"
)

func NewOpRunKilledEvent(
correlationId string,
opRunId string,
rootOpRunId string,
timestamp time.Time,
) OpRunKilledEvent {

  return OpRunKilledEvent{
    correlationId:correlationId,
    opRunId:opRunId,
    rootOpRunId:rootOpRunId,
    timestamp:timestamp,
  }

}

type OpRunKilledEvent struct {
  correlationId string
  opRunId       string
  rootOpRunId   string
  timestamp     time.Time
}

func (this OpRunKilledEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    CorrelationId string `json:"correlationId"`
    OpRunId       string `json:"opRunId"`
    RootOpRunId   string `json:"rootOpRunId"`
    Timestamp     time.Time `json:"timestamp"`
  }{
    CorrelationId:this.CorrelationId(),
    OpRunId:this.OpRunId(),
    RootOpRunId:this.RootOpRunId(),
    Timestamp:this.Timestamp(),
  }

  return json.Marshal(data)
}

func (this OpRunKilledEvent) CorrelationId() string {
  return this.correlationId
}

func (this OpRunKilledEvent) OpRunId() string {
  return this.opRunId
}

func (this OpRunKilledEvent) RootOpRunId() string {
  return this.rootOpRunId
}

func (this OpRunKilledEvent) Timestamp() time.Time {
  return this.timestamp
}
