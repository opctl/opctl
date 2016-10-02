package models

const (
  OpRunOutcomeSucceeded = "SUCCEEDED"
  OpRunOutcomeFailed = "FAILED"
  OpRunOutcomeKilled = "KILLED"
)

type OpRunEndedEvent struct {
  Descriptor *OpRunDescriptor `json:"descriptor"`
  Outcome    string `json:"outcome"`
}

