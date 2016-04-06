package dockercompose

import (
  "github.com/dev-op-spec/engine/core/models"
  "time"
  "strings"
)

func newStdErrLogWriter(
logChannel chan *models.LogEntry,
) *stdErrLogWriter {

  return &stdErrLogWriter{
    logChannel:logChannel,
  }

}

type stdErrLogWriter struct {
  logChannel chan *models.LogEntry
}

func (this stdErrLogWriter) Write(
p []byte,
) (n int, err error) {

  n = len(p)

  this.logChannel <- models.NewLogEntry(strings.TrimSpace(string(p)), time.Now().Unix(), models.StdErrStream)

  return

}
