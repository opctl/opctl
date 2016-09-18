package dockercompose

import (
  "time"
  "github.com/opspec-io/engine/core"
  "io"
  "github.com/opspec-io/sdk-golang/models"
)

func NewStdErrWriter(
eventPublisher core.EventPublisher,
opRunId string,
rootOpRunId string,
) io.Writer {

  return &stdErrWriter{
    eventPublisher:eventPublisher,
    opRunId:opRunId,
    rootOpRunId:rootOpRunId,
  }

}

// stdErrWriter implements the io.Writer interface and publishes
// written data as ContainerStdErrWrittenToEvents via an EventPublisher
type stdErrWriter struct {
  eventPublisher core.EventPublisher
  opRunId        string
  rootOpRunId    string
}

func (this stdErrWriter) Write(
p []byte,
) (n int, err error) {

  n = len(p)

  this.eventPublisher.Publish(
    models.Event{
      Timestamp:time.Now().UTC(),
      ContainerStdErrWrittenTo:&models.ContainerStdErrWrittenToEvent{
        Data:p,
        OpRunId:this.opRunId,
        RootOpRunId:this.rootOpRunId,
      },
    },
  )

  return

}
