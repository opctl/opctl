package models

import (
  "time"
  "encoding/json"
)

func NewOpRunStartedEvent(
opRunStartedAt *time.Time,
opRunId string,
) *OpRunStartedEvent {

  return &OpRunStartedEvent{
    OpRunStartedAt:opRunStartedAt,
    OpRunId:opRunId,
  }

}

type OpRunStartedEvent struct {
  OpRunStartedAt *time.Time
  OpRunId        string
}

func (this OpRunStartedEvent) Type() string {

  return "OpRunStarted"

}

func (this OpRunStartedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    OpRunStartedAt *time.Time `json:"opRunStartedAt"`
    OpRunId        string `json:"opRunId"`
    EventType      string `json:"type"`
  }{
    OpRunStartedAt:this.OpRunStartedAt,
    OpRunId:this.OpRunId,
    EventType:this.Type(),
  }

  return json.Marshal(data)
}
