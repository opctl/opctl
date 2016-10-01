package docker

import (
  "time"
  "github.com/opspec-io/engine/core"
  "io"
  "github.com/opspec-io/sdk-golang/models"
  "bufio"
)

func NewStdErrWriter(
eventPublisher core.EventPublisher,
opRunId string,
rootOpRunId string,
) io.Writer {

  reader, writer := io.Pipe()
  scanner := bufio.NewScanner(reader)

  go func() {
    for scanner.Scan() {
      eventPublisher.Publish(
        models.Event{
          Timestamp:time.Now().UTC(),
          ContainerStdErrWrittenTo:&models.ContainerStdErrWrittenToEvent{
            Data:scanner.Bytes(),
            OpRunId:opRunId,
            RootOpRunId:rootOpRunId,
          },
        },
      )
    }
  }()

  return writer

}
