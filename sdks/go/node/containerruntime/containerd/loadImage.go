package containerd

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/image/v5/copy"
	"github.com/containers/image/v5/oci/archive"
	"github.com/containers/image/v5/oci/layout"
	"github.com/containers/image/v5/signature"
	"github.com/opctl/opctl/sdks/go/model"
)

// loadImage imports an op-provided image (req.Image.Src, an OCI layout dir) into
// containerd and points req.Image.Ref at the loaded tag. It is the containerd
// analogue of the docker backend's pushImage: rather than pushing to a daemon,
// it copies the OCI layout into an OCI archive and `nerdctl load`s it.
func (cr _containerRuntime) loadImage(
	ctx context.Context,
	req *model.ContainerCall,
) error {
	imageRef := fmt.Sprintf("%s:latest", req.ContainerID)
	req.Image.Ref = &imageRef

	tmp, err := os.CreateTemp("", "opctl-image-*.tar")
	if err != nil {
		return fmt.Errorf("error loading image: %w", err)
	}
	tmpPath := tmp.Name()
	tmp.Close()
	defer os.Remove(tmpPath)

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

	srcImageRef, err := layout.NewReference(*req.Image.Src.Dir, "")
	if err != nil {
		return fmt.Errorf("error loading image: %w", err)
	}

	// oci-archive reference format is "path:tag"; the temp path has no colon, so
	// the tag (imageRef) is preserved intact when nerdctl loads the archive.
	dstImageRef, err := archive.ParseReference(fmt.Sprintf("%s:%s", tmpPath, imageRef))
	if err != nil {
		return fmt.Errorf("error loading image: %w", err)
	}

	if _, err := copy.Image(ctx, policyCtx, dstImageRef, srcImageRef, nil); err != nil {
		return fmt.Errorf("error loading image: %w", err)
	}

	if out, err := cr.nerdctl(ctx, "load", "--input", tmpPath); err != nil {
		return fmt.Errorf("error loading image: %w, %s", err, out)
	}

	return nil
}
