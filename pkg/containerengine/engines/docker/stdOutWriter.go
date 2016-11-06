package docker

import (
  "time"
  "io"
  "github.com/opspec-io/sdk-golang/pkg/models"
  "bufio"
  "github.com/opspec-io/engine/util/eventing"
)

func NewStdOutWriter(
eventPublisher eventing.EventPublisher,
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
