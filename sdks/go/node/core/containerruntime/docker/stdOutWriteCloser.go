package docker

import (
	"bufio"
	"io"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

func NewStdOutWriteCloser(
	eventPublisher pubsub.EventPublisher,
	containerID string,
	rootCallID string,
) io.WriteCloser {
	pr, pw := io.Pipe()
	go func() {
		// support lines up to 4MB
		reader := bufio.NewReaderSize(pr, 4e6)
		var err error
		var b []byte

		for {
			// chunk on newlines
			if b, err = reader.ReadBytes('\n'); len(b) > 0 {
				// always publish if len(bytes) read to ensure full stream sent; even under error conditions
				eventPublisher.Publish(
					model.Event{
						Timestamp: time.Now().UTC(),
						ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{
							Data:        b,
							ContainerID: containerID,
							RootCallID:  rootCallID,
						},
					},
				)
			}

			if err != nil {
				return
			}
		}
	}()

	return pw
}
