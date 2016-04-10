package logging

import (
  "github.com/dev-op-spec/engine/core/models"
  "time"
  "strings"
)

type LoggableIoWriter interface {
  Write(
  p []byte,
  ) (n int, err error)
}

func NewLoggableIoWriter(
logEntryOutputStream string,
logger Logger,
) LoggableIoWriter {

  return &logEmittingIoWriter{
    logEntryOutputStream:logEntryOutputStream,
    logger:logger,
  }

}

type logEmittingIoWriter struct {
  logEntryOutputStream string
  logger               Logger
}

func (this logEmittingIoWriter) Write(
p []byte,
) (n int, err error) {

  n = len(p)

  this.logger(
    models.NewLogEntryEmittedEvent(
      time.Now(),
      strings.TrimSpace(string(p)),
      this.logEntryOutputStream,
    ),
  )

  return

}
