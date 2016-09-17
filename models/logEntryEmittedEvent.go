package models

type LogEntryEmittedEvent struct {
  Msg          string `json:"logEntryMsg"`
  OutputStream string `json:"logEntryOutputStream"`
}
