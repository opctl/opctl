package docker

import (
	"context"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker/daemon"
	"github.com/containers/image/v5/oci/layout"
	"github.com/containers/image/v5/signature"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/pkg/errors"
)

//counterfeiter:generate -o internal/fakes/imagePusher.go . imagePusher
type imagePusher interface {
	Push(
		ctx context.Context,
		imageRef string,
		imageSrc *model.Value,
	) error
}

func newImagePusher() imagePusher {
	return _imagePusher{}
}

type _imagePusher struct{}

func (ip _imagePusher) Push(
	ctx context.Context,
	imageRef string,
	imageSrc *model.Value,
) error {
	policyCtx, err := signature.NewPolicyContext(
		&signature.Policy{
			Default: []signature.PolicyRequirement{
				signature.NewPRInsecureAcceptAnything(),
			},
		},
	)
	if err != nil {
		return errors.Wrap(err, "error loading image")
	}

	srcImageRef, err := layout.NewReference(*imageSrc.Dir, "")
	if err != nil {
		return errors.Wrap(err, "error loading image")
	}

	dstImageRef, err := daemon.ParseReference(imageRef)
	if err != nil {
		return errors.Wrap(err, "error loading image")
	}

	if _, err := copy.Image(ctx, policyCtx, dstImageRef, srcImageRef, nil); err != nil {
		return errors.Wrap(err, "error loading image")
	}

	return nil
}
