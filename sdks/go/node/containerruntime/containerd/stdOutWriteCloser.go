package containerd

import (
	"bufio"
	"io"
	"time"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

// newStdOutWriteCloser returns an io.WriteCloser that republishes everything
// written to it as ContainerStdOutWrittenTo events (chunked on newlines). Used
// to surface `nerdctl pull` progress through opctl's event stream, mirroring
// the docker backend's writer.
func newStdOutWriteCloser(
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
