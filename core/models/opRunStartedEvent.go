package models

import (
  "time"
  "encoding/json"
)

func NewOpRunStartedEvent(
correlationId string,
timestamp time.Time,
parentOpRunId string,
opRunOpUrl Url,
opRunId string,
) OpRunStartedEvent {

  return &opRunStartedEvent{
    correlationId:correlationId,
    opRunOpUrl:opRunOpUrl,
    opRunId:opRunId,
    parentOpRunId:parentOpRunId,
    timestamp:timestamp,
  }

}

type OpRunStartedEvent interface {
  CorrelationId() string
  OpRunId() string
  OpRunOpUrl() Url
  ParentOpRunId() string
  Timestamp() time.Time
}

type opRunStartedEvent struct {
  correlationId string
  opRunId       string
  opRunOpUrl    Url
  parentOpRunId string
  timestamp     time.Time
}

func (this opRunStartedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    CorrelationId string `json:"correlationId"`
    OpRunId       string `json:"opRunId"`
    OpRunOpUrl    string `json:"opRunOpUrl"`
    ParentOpRunId string `json:"parentOpRunId"`
    Timestamp     time.Time `json:"timestamp"`
  }{
    CorrelationId:this.CorrelationId(),
    OpRunId:this.opRunId,
    OpRunOpUrl:this.opRunOpUrl.String(),
    ParentOpRunId:this.parentOpRunId,
    Timestamp:this.timestamp,
  }

  return json.Marshal(data)
}

func (this opRunStartedEvent) CorrelationId() string {
  return this.correlationId
}

func (this opRunStartedEvent) OpRunId() string {
  return this.opRunId
}

func (this opRunStartedEvent) OpRunOpUrl() Url {
  return this.opRunOpUrl
}

func (this opRunStartedEvent) ParentOpRunId() string {
  return this.parentOpRunId
}

func (this opRunStartedEvent) Timestamp() time.Time {
  return this.timestamp
}
