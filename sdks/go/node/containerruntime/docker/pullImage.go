package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/distribution/reference"
	"github.com/docker/docker/api/types/image"
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

func pullImage(
	ctx context.Context,
	containerCall *model.ContainerCall,
	dockerClient dockerClientPkg.CommonAPIClient,
	rootCallID string,
	eventPublisher pubsub.EventPublisher,
) error {
	imageRef := *containerCall.Image.Ref

	needsPull, err := doesImageNeedPull(ctx, dockerClient, imageRef, eventPublisher)
	if err != nil {
		return err
	}
	if !needsPull {
		eventPublisher.Publish(model.Event{
			Timestamp: time.Now().UTC(),
			ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{
				Data:        []byte(fmt.Sprintf("Skipping image pull: %s\n", imageRef)),
				OpRef:       containerCall.OpPath,
				ContainerID: containerCall.ContainerID,
				RootCallID:  rootCallID,
			},
		})
		return nil
	}

	imagePullCreds := containerCall.Image.PullCreds
	containerID := containerCall.ContainerID

	platformParts := []string{
		"linux",
	}

	if containerCall.Image.Platform != nil && containerCall.Image.Platform.Arch != nil {
		platformParts = append(platformParts, *containerCall.Image.Platform.Arch)
	}

	imagePullOptions := image.PullOptions{
		Platform: strings.Join(
			platformParts,
			"/",
		),
	}

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

	imagePullResp, err := dockerClient.ImagePull(
		ctx,
		imageRef,
		imagePullOptions,
	)
	if imagePullResp == nil || err != nil {
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

func doesImageNeedPull(
	ctx context.Context,
	dockerClient dockerClientPkg.CommonAPIClient,
	imageRef string,
	eventPublisher pubsub.EventPublisher,
) (bool, error) {
	// Skip pulling for non-tagged images that already are present
	// This reduces the chance of hitting docker rate limiting errors
	// and speeds up execution.
	ref, err := reference.ParseAnyReference(strings.ToLower(imageRef))
	if err != nil {
		return true, err
	}
	named, err := reference.ParseNormalizedNamed(ref.String())
	if err != nil {
		return true, err
	}
	if tagged, ok := named.(reference.Tagged); ok && tagged.Tag() != "latest" {
		_, _, err := dockerClient.ImageInspectWithRaw(ctx, imageRef)
		if err == nil {
			return false, nil
		}
		// this err can be ignored, since it's expected to be "image not found"
	}

	return true, nil
}
