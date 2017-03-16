package docker

import (
	"bufio"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
	"io"
	"time"
)

func NewStdOutWriter(
	eventPublisher pubsub.EventPublisher,
	containerId string,
	rootOpId string,
) io.Writer {

	reader, writer := io.Pipe()
	scanner := bufio.NewScanner(reader)

	go func() {
		for scanner.Scan() {
			eventPublisher.Publish(
				&model.Event{
					Timestamp: time.Now().UTC(),
					ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenToEvent{
						Data:        scanner.Bytes(),
						ContainerId: containerId,
						RootOpId:    rootOpId,
					},
				},
			)
		}
	}()

	return writer

}
