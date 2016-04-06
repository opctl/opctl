package dockercompose

import (
  "github.com/dev-op-spec/engine/core/models"
  "time"
  "strings"
)

func newStdOutLogWriter(
logChannel chan *models.LogEntry,
) *stdOutLogWriter {

  return &stdOutLogWriter{
    logChannel:logChannel,
  }

}

type stdOutLogWriter struct {
  logChannel chan *models.LogEntry
}

func (this stdOutLogWriter) Write(
p []byte,
) (n int, err error) {

  n = len(p)

  this.logChannel <- models.NewLogEntry(strings.TrimSpace(string(p)), time.Now().Unix(), models.StdOutStream)

  return

}
