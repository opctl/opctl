package models

import (
  "time"
  "encoding/json"
)

func NewOpRunFinishedEvent(
opRunFinishTime time.Time,
opRunExitCode int,
opRunId string,
) OpRunFinishedEvent {

  return &opRunFinishedEvent{
    opRunFinishTime:opRunFinishTime,
    opRunExitCode:opRunExitCode,
    opRunId:opRunId,
  }

}

type OpRunFinishedEvent interface {
  OpRunFinishTime() time.Time
  OpRunExitCode() int
  OpRunId() string
  Type() string
}

type opRunFinishedEvent struct {
  opRunFinishTime time.Time
  opRunExitCode   int
  opRunId         string
}

func (this opRunFinishedEvent) MarshalJSON() ([]byte, error) {

  data := struct {
    OpRunFinishTime time.Time `json:"opRunFinishTime"`
    OpRunExitCode   int `json:"opRunExitCode"`
    OpRunId         string `json:"opRunId"`
    Type            string `json:"type"`
  }{
    OpRunFinishTime:this.OpRunFinishTime(),
    OpRunExitCode:this.OpRunExitCode(),
    OpRunId:this.OpRunId(),
    Type:this.Type(),
  }

  return json.Marshal(data)
}

func (this opRunFinishedEvent) OpRunFinishTime() time.Time {
  return this.opRunFinishTime
}

func (this opRunFinishedEvent) OpRunExitCode() int {
  return this.opRunExitCode
}

func (this opRunFinishedEvent) OpRunId() string {
  return this.opRunId
}

func (this opRunFinishedEvent) Type() string {

  return "OpRunFinished"

}

