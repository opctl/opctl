package models

import (
  "time"
  "encoding/json"
)

func NewOpRunStartedEvent(
correlationId string,
opRunOpUrl string,
opRunId string,
parentOpRunId string,
rootOpRunId string,
timestamp time.Time,
) OpRunStartedEvent {

  return OpRunStartedEvent{
    correlationId:correlationId,
    opRunOpUrl:opRunOpUrl,
    opRunId:opRunId,
    parentOpRunId:parentOpRunId,
    rootOpRunId:rootOpRunId,
    timestamp:timestamp,
  }

}

type OpRunStartedEvent struct {
  correlationId string
  opRunId       string
  opRunOpUrl    string
  parentOpRunId string
  rootOpRunId   string
  timestamp     time.Time
}

func (this OpRunStartedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    CorrelationId string `json:"correlationId"`
    OpRunId       string `json:"opRunId"`
    OpRunOpUrl    string `json:"opRunOpUrl"`
    ParentOpRunId string `json:"parentOpRunId"`
    RootOpRunId   string `json:"rootOpRunId"`
    Timestamp     time.Time `json:"timestamp"`
  }{
    CorrelationId:this.CorrelationId(),
    OpRunId:this.OpRunId(),
    OpRunOpUrl:this.OpRunOpUrl(),
    ParentOpRunId:this.ParentOpRunId(),
    RootOpRunId:this.RootOpRunId(),
    Timestamp:this.Timestamp(),
  }

  return json.Marshal(data)
}

func (this OpRunStartedEvent) CorrelationId() string {
  return this.correlationId
}

func (this OpRunStartedEvent) OpRunId() string {
  return this.opRunId
}

func (this OpRunStartedEvent) OpRunOpUrl() string {
  return this.opRunOpUrl
}

func (this OpRunStartedEvent) ParentOpRunId() string {
  return this.parentOpRunId
}

func (this OpRunStartedEvent) RootOpRunId() string {
  return this.rootOpRunId
}

func (this OpRunStartedEvent) Timestamp() time.Time {
  return this.timestamp
}
