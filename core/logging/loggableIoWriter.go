package logging

import (
  "github.com/open-devops/engine/core/models"
  "time"
  "strings"
)

type LoggableIoWriter interface {
  Write(
  p []byte,
  ) (n int, err error)
}

func NewLoggableIoWriter(
correlationId string,
logEntryOutputStream string,
logger Logger,
) LoggableIoWriter {

  return &logEmittingIoWriter{
    correlationId:correlationId,
    logEntryOutputStream:logEntryOutputStream,
    logger:logger,
  }

}

type logEmittingIoWriter struct {
  correlationId        string
  logEntryOutputStream string
  logger               Logger
}

func (this logEmittingIoWriter) Write(
p []byte,
) (n int, err error) {

  n = len(p)

  this.logger(
    models.NewLogEntryEmittedEvent(
      this.correlationId,
      time.Now().UTC(),
      strings.TrimSpace(string(p)),
      this.logEntryOutputStream,
    ),
  )

  return

}
