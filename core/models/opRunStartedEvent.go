package models

import (
  "time"
  "encoding/json"
)

func NewOpRunStartedEvent(
opRunStartTime *time.Time,
opRunId string,
) *OpRunStartedEvent {

  return &OpRunStartedEvent{
    OpRunStartTime:opRunStartTime,
    OpRunId:opRunId,
  }

}

type OpRunStartedEvent struct {
  OpRunStartTime *time.Time
  OpRunId        string
}

func (this OpRunStartedEvent) Type() string {

  return "OpRunStarted"

}

func (this OpRunStartedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    OpRunStartTime *time.Time `json:"opRunStartTime"`
    OpRunId        string `json:"opRunId"`
    EventType      string `json:"type"`
  }{
    OpRunStartTime:this.OpRunStartTime,
    OpRunId:this.OpRunId,
    EventType:this.Type(),
  }

  return json.Marshal(data)
}
