package models

import "time"

func NewLogEntry(
message string,
timestamp time.Time,
stream string,
) *LogEntry {

  return &LogEntry{
    Message:message,
    Timestamp:timestamp,
    Stream:stream,
  }

}

type LogEntry struct {
  Message   string `json:"message"`
  Timestamp time.Time `json:"timestamp"`
  Stream    string `json:"stream"`
}

const StdOutStream string = "stdout"

const StdErrStream string = "stderr"
