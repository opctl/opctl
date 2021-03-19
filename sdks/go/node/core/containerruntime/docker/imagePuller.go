package docker

import (
	"context"
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types"
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o internal/fakes/imagePuller.go . imagePuller
type imagePuller interface {
	Pull(
		ctx context.Context,
		containerID string,
		imagePullCreds *model.Creds,
		imageRef string,
		rootCallID string,
		eventPublisher pubsub.EventPublisher,
	) error
}

func newImagePuller(
	dockerClient dockerClientPkg.CommonAPIClient,
) imagePuller {
	return _imagePuller{
		dockerClient,
	}
}

type _imagePuller struct {
	dockerClient dockerClientPkg.CommonAPIClient
}

func (ip _imagePuller) Pull(
	ctx context.Context,
	containerID string,
	imagePullCreds *model.Creds,
	imageRef string,
	rootCallID string,
	eventPublisher pubsub.EventPublisher,
) error {

	imagePullOptions := types.ImagePullOptions{}
	if imagePullCreds != nil &&
		imagePullCreds.Username != "" &&
		imagePullCreds.Password != "" {
		var err error
		imagePullOptions.RegistryAuth, err = constructRegistryAuth(
			imagePullCreds.Username,
			imagePullCreds.Password,
		)
		if err != nil {
			return err
		}
	}

	imagePullResp, err := ip.dockerClient.ImagePull(
		ctx,
		imageRef,
		imagePullOptions,
	)
	if err != nil {
		return err
	}
	defer imagePullResp.Close()

	stdOutWriter := NewStdOutWriteCloser(eventPublisher, containerID, rootCallID)
	defer stdOutWriter.Close()

	dec := json.NewDecoder(imagePullResp)
	for {
		var jm jsonmessage.JSONMessage
		if err = dec.Decode(&jm); err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}
		jm.Display(stdOutWriter, false)
	}
}
