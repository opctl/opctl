package models

import (
  "time"
  "encoding/json"
)

func NewOpRunStartedEvent(
timestamp time.Time,
opRunParentId *string,
opRunOpUrl Url,
opRunId string,
) OpRunStartedEvent {

  return &opRunStartedEvent{
    timestamp:timestamp,
    opRunParentId:opRunParentId,
    opRunOpUrl:opRunOpUrl,
    opRunId:opRunId,
  }

}

type OpRunStartedEvent interface {
  OpRunId() string
  OpRunOpUrl() Url
  OpRunParentId() *string
  Timestamp() time.Time
  Type() string
}

type opRunStartedEvent struct {
  opRunId       string
  opRunOpUrl    Url
  opRunParentId *string
  timestamp     time.Time
}

func (this opRunStartedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    OpRunId       string `json:"opRunId"`
    OpRunOpUrl    string `json:"opRunOpUrl"`
    OpRunParentId *string `json:"opRunParentId"`
    Timestamp     time.Time `json:"timestamp"`
    Type          string `json:"type"`
  }{
    OpRunId:this.opRunId,
    OpRunOpUrl:this.opRunOpUrl.String(),
    OpRunParentId:this.opRunParentId,
    Timestamp:this.timestamp,
    Type:this.Type(),
  }

  return json.Marshal(data)
}

func (this opRunStartedEvent) OpRunId() string {
  return this.opRunId
}

func (this opRunStartedEvent) OpRunOpUrl() Url {
  return this.opRunOpUrl
}

func (this opRunStartedEvent) OpRunParentId() *string {
  return this.opRunParentId
}

func (this opRunStartedEvent) Timestamp() time.Time {
  return this.timestamp
}

func (this opRunStartedEvent) Type() string {

  return "OpRunStarted"

}
