package docker

import (
	"bufio"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/util/pubsub"
	"io"
	"time"
)

func NewStdOutWriteCloser(
	eventPublisher pubsub.EventPublisher,
	containerID string,
	rootOpID string,
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
						ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenToEvent{
							Data:        b,
							ContainerID: containerID,
							RootOpID:    rootOpID,
						},
					},
				)
			}

			if nil != err {
				return
			}
		}
	}()

	return pw

}
