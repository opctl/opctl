package models

type OpRunStartedEvent struct {
  OpRunId       string `json:"opRunId"`
  OpRunOpUrl    string `json:"opRunOpUrl"`
  ParentOpRunId string `json:"parentOpRunId"`
  RootOpRunId   string `json:"rootOpRunId"`
}
