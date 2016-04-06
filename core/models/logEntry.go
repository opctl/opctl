package models

func NewLogEntry(
message string,
timestamp int64,
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
  Timestamp int64 `json:"timestamp"`
  Stream    string `json:"stream"`
}
