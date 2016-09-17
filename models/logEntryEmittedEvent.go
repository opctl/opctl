package models

type LogEntryEmittedEvent struct {
  LogEntryMsg          string `json:"logEntryMsg"`
  LogEntryOutputStream string `json:"logEntryOutputStream"`
}
