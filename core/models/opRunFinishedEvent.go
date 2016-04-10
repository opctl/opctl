package models

import (
  "time"
  "encoding/json"
)

func NewOpRunFinishedEvent(
opRunFinishedAt *time.Time,
opRunId string,
) *OpRunFinishedEvent {

  return &OpRunFinishedEvent{
    OpRunFinishedAt:opRunFinishedAt,
    OpRunId:opRunId,
  }

}

type OpRunFinishedEvent struct {
  OpRunFinishedAt *time.Time
  OpRunId         string
}

func (this OpRunFinishedEvent) Type() string {

  return "OpRunFinished"

}

func (this OpRunFinishedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    OpRunFinishedAt *time.Time `json:"opRunFinishedAt"`
    OpRunId         string `json:"opRunId"`
    EventType       string `json:"type"`
  }{
    OpRunFinishedAt:this.OpRunFinishedAt,
    OpRunId:this.OpRunId,
    EventType:this.Type(),
  }

  return json.Marshal(data)
}