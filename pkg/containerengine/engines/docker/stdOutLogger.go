package docker

import (
	"bufio"
	"github.com/docker/docker/api/types"
	"github.com/opspec-io/opctl/util/eventbus"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"golang.org/x/net/context"
	"io"
	"time"
)

func (this _containerEngine) stdOutLogger(
	eventPublisher eventbus.EventPublisher,
	containerId string,
	opGraphId string,
) (err error) {

	var readCloser io.ReadCloser
	readCloser, err = this.dockerEngine.ContainerLogs(
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
	scanner := bufio.NewScanner(readCloser)

	go func() {
		defer readCloser.Close()
		for scanner.Scan() {
			eventPublisher.Publish(
				model.Event{
					Timestamp: time.Now().UTC(),
					ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenToEvent{
						Data:        scanner.Bytes(),
						ContainerId: containerId,
						OpGraphId:   opGraphId,
					},
				},
			)
		}
	}()
	return
}
