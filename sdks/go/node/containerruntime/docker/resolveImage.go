package docker

import (
	"context"
	"fmt"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

func (cr _containerRuntime) ResolveImage(ctx context.Context, image *model.ContainerCallImage, root string, events pubsub.EventPublisher) (*containerruntime.ImageIdentity, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if image == nil || image.Ref == nil {
		return nil, fmt.Errorf("Go bootstrap needs an image reference")
	}
	// A failed pull may still leave a usable local image for offline execution.
	_ = pullImage(ctx, &model.ContainerCall{Image: image}, cr.dockerClient, root, events)
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	info, _, err := cr.dockerClient.ImageInspectWithRaw(ctx, *image.Ref)
	if err != nil || info.ID == "" || info.Os != "linux" || info.Architecture == "" {
		return nil, fmt.Errorf("cannot resolve immutable Linux Go bootstrap image %q", *image.Ref)
	}
	return &containerruntime.ImageIdentity{Ref: info.ID, ID: info.ID, OS: info.Os, Arch: info.Architecture, Variant: info.Variant}, nil
}
