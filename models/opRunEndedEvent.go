package models

const (
  OpRunOutcomeSucceeded = "SUCCEEDED"
  OpRunOutcomeFailed = "FAILED"
  OpRunOutcomeKilled = "KILLED"
)

type OpRunEndedEvent struct {
  OpRunId       string `json:"opRunId"`
  Outcome       string `json:"outcome"`
  ParentOpRunId string `json:"parentOpRunId"`
  RootOpRunId   string `json:"rootOpRunId"`
}

