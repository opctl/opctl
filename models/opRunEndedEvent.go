package models

type OpRunEndedEvent struct {
  Id      string `json:"id"`
  Outcome string `json:"outcome"`
  RootId  string `json:"rootId"`
}

