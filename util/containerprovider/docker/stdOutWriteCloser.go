package docker

import (
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
	"io"
	"time"
)

func NewStdOutWriteCloser(
	eventPublisher pubsub.EventPublisher,
	containerId string,
	rootOpId string,
) io.WriteCloser {
	reader, writer := io.Pipe()
	go func() {
		chunk := make([]byte, 1024)
		var n int
		var err error

		for {
			// rather than chunking by line, we chunk by time at a rate of 30 FPS (frames per second)
			// why? chunking by line would make TTY behaviors such as line editing behave non-TTY like
			<-time.After(33 * time.Millisecond)
			n, err = reader.Read(chunk)
			if n > 0 {
				// always call onChunk if n > 0 to ensure full stream sent; even under error conditions
				eventPublisher.Publish(
					&model.Event{
						Timestamp: time.Now().UTC(),
						ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenToEvent{
							Data:        chunk[0:n],
							ContainerId: containerId,
							RootOpId:    rootOpId,
						},
					},
				)
			}
			if nil != err {
				return
			}
		}
	}()

	return writer

}
