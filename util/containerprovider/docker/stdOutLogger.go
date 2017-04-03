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

func (this _containerProvider) stdOutLogger(
	eventPublisher pubsub.EventPublisher,
	containerId string,
	imageRef string,
	pkgRef string,
	rootOpId string,
) (err error) {

	var readCloser io.ReadCloser
	readCloser, err = this.dockerClient.ContainerLogs(
		context.Background(),
		containerId,
		types.ContainerLogsOptions{
			Follow:     true,
			ShowStdout: true,
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
					ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenToEvent{
						Data:         scanner.Bytes(),
						ContainerId:  containerId,
						ContainerRef: imageRef,
						PkgRef:       pkgRef,
						RootOpId:     rootOpId,
					},
				},
			)
		}
		defer readCloser.Close()
	}()
	return
}
