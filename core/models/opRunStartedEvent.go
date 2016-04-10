package models

import (
  "time"
  "encoding/json"
)

func NewOpRunStartedEvent(
opRunStartTime time.Time,
opRunParentId *string,
opRunOpUrl Url,
opRunId string,
) *OpRunStartedEvent {

  return &OpRunStartedEvent{
    OpRunStartTime:opRunStartTime,
    OpRunParentId:opRunParentId,
    OpRunOpUrl:opRunOpUrl,
    OpRunId:opRunId,
  }

}

type OpRunStartedEvent struct {
  OpRunStartTime time.Time
  OpRunParentId  *string
  OpRunOpUrl     Url
  OpRunId        string
}

func (this OpRunStartedEvent) Type() string {

  return "OpRunStarted"

}

func (this OpRunStartedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    OpRunStartTime time.Time `json:"opRunStartTime"`
    OpRunParentId  *string `json:"opRunParentId"`
    OpRunOpUrl     string `json:"opRunOpUrl"`
    OpRunId        string `json:"opRunId"`
    Type           string `json:"type"`
  }{
    OpRunStartTime:this.OpRunStartTime,
    OpRunParentId:this.OpRunParentId,
    OpRunOpUrl:this.OpRunOpUrl.String(),
    OpRunId:this.OpRunId,
    Type:this.Type(),
  }

  return json.Marshal(data)
}
