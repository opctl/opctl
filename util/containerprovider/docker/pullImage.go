package docker

import (
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/reference"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/opctl/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
	"golang.org/x/net/context"
	"io"
)

func (this _containerProvider) pullImage(
	dcgContainerImage *model.DCGContainerCallImage,
	containerId string,
	rootOpId string,
	eventPublisher pubsub.EventPublisher,
) (err error) {
	// ensure tag present in image string.
	// if not present, docker defaults to downloading all tags
	imageName, tag, err := reference.Parse(dcgContainerImage.Ref)
	if nil != err {
		return
	}
	imageRef := fmt.Sprintf("%v:%v", imageName, tag)

	imagePullOptions := types.ImagePullOptions{}
	if nil != dcgContainerImage.PullAuth &&
		"" != dcgContainerImage.PullAuth.Username &&
		"" != dcgContainerImage.PullAuth.Password {
		imagePullOptions.RegistryAuth, err = constructRegistryAuth(
			dcgContainerImage.PullAuth.Username,
			dcgContainerImage.PullAuth.Password,
		)
		fmt.Printf("imagePullOptions.RegistryAuth: %v \n", imagePullOptions.RegistryAuth)
		if nil != err {
			return
		}
	}

	imagePullResp, err := this.dockerClient.ImagePull(
		context.Background(),
		imageRef,
		imagePullOptions,
	)
	if nil != err {
		return
	}

	defer imagePullResp.Close()

	stdOutWriter := NewStdOutWriter(eventPublisher, containerId, rootOpId)
	dec := json.NewDecoder(imagePullResp)
	for {
		var jm jsonmessage.JSONMessage
		if err = dec.Decode(&jm); nil != err {
			if err == io.EOF {
				err = nil
			}
			return
		}
		jm.Display(stdOutWriter, false)
	}
}
