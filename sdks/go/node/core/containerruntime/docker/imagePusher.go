package docker

import (
	"context"
	"fmt"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker/daemon"
	"github.com/containers/image/v5/oci/layout"
	"github.com/containers/image/v5/signature"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o internal/fakes/imagePusher.go . imagePusher
type imagePusher interface {
	Push(
		ctx context.Context,
		containerID string,
		imageRef string,
		imageSrc *model.Value,
		rootOpID string,
		eventPublisher pubsub.EventPublisher,
	) error
}

func newImagePusher() imagePusher {
	return _imagePusher{}
}

type _imagePusher struct{}

func (ip _imagePusher) Push(
	ctx context.Context,
	containerID string,
	imageRef string,
	imageSrc *model.Value,
	rootOpID string,
	eventPublisher pubsub.EventPublisher,
) error {

	policyCtx, policyCtxErr := signature.NewPolicyContext(
		&signature.Policy{Default: []signature.PolicyRequirement{signature.NewPRInsecureAcceptAnything()}},
	)
	if nil != policyCtxErr {
		return fmt.Errorf("error encountered loading image; error was: %v", policyCtxErr)
	}

	srcImageRef, srcErr := layout.NewReference(
		*imageSrc.Dir,
		"",
	)
	if nil != srcErr {
		return fmt.Errorf("error encountered loading image; error was: %v", srcErr)
	}

	dstImageRef, dstErr := daemon.ParseReference(
		imageRef,
	)
	if nil != dstErr {
		return fmt.Errorf("error encountered loading image; error was: %v", dstErr)
	}

	_, copyErr := copy.Image(
		ctx,
		policyCtx,
		dstImageRef,
		srcImageRef,
		nil,
	)
	if nil != copyErr {
		return fmt.Errorf("error encountered loading image; error was: %v", copyErr)
	}

	return nil
}
