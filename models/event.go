package models

import "time"

type Event struct {
  Timestamp       time.Time `json:"timestamp"`
  OpRunEnded      *OpRunEndedEvent `yaml:"opRunEnded,omitempty"`
  OpRunStarted    *OpRunStartedEvent `yaml:"opRunStarted,omitempty"`
  LogEntryEmitted *LogEntryEmittedEvent `yaml:"logEntryEmitted,omitempty"`
}
