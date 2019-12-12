package docker

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeImagePusher.go --fake-name fakeImagePusher ./ imagePusher

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

type imagePusher interface {
	Push(
		ctx context.Context,
		dcgContainerImage *model.DCGContainerCallImage,
		containerID string,
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
	dcgContainerImage *model.DCGContainerCallImage,
	containerID string,
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
		*dcgContainerImage.Src.Dir,
		"",
	)
	if nil != srcErr {
		return fmt.Errorf("error encountered loading image; error was: %v", srcErr)
	}

	dstImageRef, dstErr := daemon.ParseReference(
		dcgContainerImage.Ref,
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
