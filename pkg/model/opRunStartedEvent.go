package models

type OpRunStartedEvent struct {
  OpRunId       string `json:"opRunId"`
  OpRef         string `json:"opRef"`
  RootOpRunId   string `json:"rootOpRunId"`
}
