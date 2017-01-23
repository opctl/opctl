package docker

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/reference"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/opspec-io/opctl/util/eventbus"
	"golang.org/x/net/context"
	"io"
)

func (this _containerEngine) pullImage(
	imageRef string,
	containerId string,
	opGraphId string,
	eventPublisher eventbus.EventPublisher,
) (err error) {
	// ensure tag present in image string.
	// if not present, docker defaults to downloading all tags
	imageName, tag, err := reference.Parse(imageRef)
	if err != nil {
		return
	}
	imageRef = fmt.Sprintf("%v:%v", imageName, tag)

	imagePullResp, err := this.dockerClient.ImagePull(
		context.Background(),
		imageRef,
		types.ImagePullOptions{},
	)
	if nil != err {
		return
	}

	defer imagePullResp.Close()

	stdOutWriter := NewStdOutWriter(eventPublisher, containerId, opGraphId)
	dec := json.NewDecoder(imagePullResp)
	for {
		var jm jsonmessage.JSONMessage
		if err = dec.Decode(&jm); err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		jm.Display(stdOutWriter, false)
	}
}
