package docker

import (
  "time"
  "io"
  "github.com/opspec-io/sdk-golang/models"
  "github.com/opspec-io/engine/core"
  "bufio"
)

func NewStdOutWriter(
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
          ContainerStdOutWrittenTo:&models.ContainerStdOutWrittenToEvent{
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
