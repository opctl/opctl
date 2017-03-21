package docker

import (
	"bufio"
	"github.com/docker/docker/api/types"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
	"golang.org/x/net/context"
	"io"
	"time"
)

func (this _containerProvider) stdErrLogger(
	eventPublisher pubsub.EventPublisher,
	containerId string,
	rootOpId string,
) (err error) {

	var readCloser io.ReadCloser
	readCloser, err = this.dockerClient.ContainerLogs(
		context.Background(),
		containerId,
		types.ContainerLogsOptions{
			Follow:     true,
			ShowStderr: true,
			Details:    false,
		},
	)
	if nil != err {
		return
	}

	go func() {
		scanner := bufio.NewScanner(readCloser)
		for scanner.Scan() {
			eventPublisher.Publish(
				&model.Event{
					Timestamp: time.Now().UTC(),
					ContainerStdErrWrittenTo: &model.ContainerStdErrWrittenToEvent{
						Data:        scanner.Bytes(),
						ContainerId: containerId,
						RootOpId:    rootOpId,
					},
				},
			)
		}
		defer readCloser.Close()
	}()
	return
}
