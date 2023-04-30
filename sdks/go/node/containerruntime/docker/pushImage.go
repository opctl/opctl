package docker

import (
	"context"
	"fmt"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/docker/daemon"
	"github.com/containers/image/v5/oci/layout"
	"github.com/containers/image/v5/signature"
	"github.com/opctl/opctl/sdks/go/model"
)

func pushImage(
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
		return fmt.Errorf("error loading image: %w", err)
	}

	srcImageRef, err := layout.NewReference(*imageSrc.Dir, "")
	if err != nil {
		return fmt.Errorf("error loading image: %w", err)
	}

	dstImageRef, err := daemon.ParseReference(imageRef)
	if err != nil {
		return fmt.Errorf("error loading image: %w", err)
	}

	if _, err := copy.Image(ctx, policyCtx, dstImageRef, srcImageRef, nil); err != nil {
		return fmt.Errorf("error loading image: %w", err)
	}

	return nil
}
