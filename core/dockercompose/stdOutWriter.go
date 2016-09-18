package dockercompose

import (
  "time"
  "io"
  "github.com/opspec-io/sdk-golang/models"
  "github.com/opspec-io/engine/core"
)

func NewStdOutWriter(
eventPublisher core.EventPublisher,
opRunId string,
rootOpRunId string,
) io.Writer {

  return &stdOutWriter{
    eventPublisher:eventPublisher,
    opRunId:opRunId,
    rootOpRunId:rootOpRunId,
  }

}

// stdOutWriter implements the io.Writer interface and publishes
// written data as ContainerStdOutWrittenToEvents via an EventPublisher
type stdOutWriter struct {
  eventPublisher core.EventPublisher
  opRunId        string
  rootOpRunId    string
}

func (this stdOutWriter) Write(
p []byte,
) (n int, err error) {

  n = len(p)

  this.eventPublisher.Publish(
    models.Event{
      Timestamp:time.Now().UTC(),
      ContainerStdOutWrittenTo:&models.ContainerStdOutWrittenToEvent{
        Data:p,
        OpRunId:this.opRunId,
        RootOpRunId:this.rootOpRunId,
      },
    },
  )

  return

}
