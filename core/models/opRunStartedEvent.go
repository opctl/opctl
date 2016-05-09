package models

import (
  "time"
  "encoding/json"
)

func NewOpRunStartedEvent(
timestamp time.Time,
parentOpRunId string,
opRunOpUrl Url,
opRunId string,
) OpRunStartedEvent {

  return &opRunStartedEvent{
    opRunOpUrl:opRunOpUrl,
    opRunId:opRunId,
    parentOpRunId:parentOpRunId,
    timestamp:timestamp,
  }

}

type OpRunStartedEvent interface {
  OpRunId() string
  OpRunOpUrl() Url
  ParentOpRunId() string
  Timestamp() time.Time
}

type opRunStartedEvent struct {
  opRunId       string
  opRunOpUrl    Url
  parentOpRunId string
  timestamp     time.Time
}

func (this opRunStartedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    OpRunId       string `json:"opRunId"`
    OpRunOpUrl    string `json:"opRunOpUrl"`
    ParentOpRunId string `json:"parentOpRunId"`
    Timestamp     time.Time `json:"timestamp"`
  }{
    OpRunId:this.opRunId,
    OpRunOpUrl:this.opRunOpUrl.String(),
    ParentOpRunId:this.parentOpRunId,
    Timestamp:this.timestamp,
  }

  return json.Marshal(data)
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
