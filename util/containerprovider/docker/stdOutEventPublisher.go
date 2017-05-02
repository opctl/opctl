package docker

import (
	"bufio"
	"github.com/docker/docker/api/types"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
	"golang.org/x/net/context"
	"time"
)

func (this _containerProvider) stdOutEventPublisher(
	eventPublisher pubsub.EventPublisher,
	containerId string,
	imageRef string,
	pkgRef string,
	rootOpId string,
) error {

	readCloser, err := this.dockerClient.ContainerLogs(
		context.Background(),
		containerId,
		types.ContainerLogsOptions{
			Follow:     true,
			ShowStdout: true,
		},
	)
	defer readCloser.Close()
	if nil != err {
		return err
	}

	scanner := bufio.NewScanner(readCloser)

	// scan writes until EOF or error
	for scanner.Scan() {

		// publish writes
		eventPublisher.Publish(
			&model.Event{
				Timestamp: time.Now().UTC(),
				ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenToEvent{
					Data:        scanner.Bytes(),
					ContainerId: containerId,
					ImageRef:    imageRef,
					PkgRef:      pkgRef,
					RootOpId:    rootOpId,
				},
			},
		)

	}

	return scanner.Err()

}
