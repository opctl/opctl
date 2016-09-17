package models

type OpRunEndedEvent struct {
  OpRunId       string `json:"opRunId"`
  Outcome       string `json:"outcome"`
  RootOpRunId   string `json:"rootOpRunId"`
}

