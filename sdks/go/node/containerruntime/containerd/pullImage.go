package containerd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/distribution/reference"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
)

// pullImage pulls the call's image via nerdctl, streaming progress to the event
// publisher. Registry auth and mirror redirection are resolved by
// containerd/nerdctl host config (`/etc/containerd/certs.d/<host>/hosts.toml`
// and `~/.docker/config.json` credential helpers) rather than per-op pullCreds.
func (cr _containerRuntime) pullImage(
	ctx context.Context,
	req *model.ContainerCall,
	rootCallID string,
	eventPublisher pubsub.EventPublisher,
) error {
	imageRef := *req.Image.Ref

	if !cr.needsPull(ctx, imageRef) {
		eventPublisher.Publish(model.Event{
			Timestamp: time.Now().UTC(),
			ContainerStdOutWrittenTo: &model.ContainerStdOutWrittenTo{
				Data:        []byte(fmt.Sprintf("Skipping image pull: %s\n", imageRef)),
				OpRef:       req.OpPath,
				ContainerID: req.ContainerID,
				RootCallID:  rootCallID,
			},
		})
		return nil
	}

	args := []string{"pull"}
	if req.Image.Platform != nil && req.Image.Platform.Arch != nil {
		args = append(args, "--platform", fmt.Sprintf("linux/%s", *req.Image.Platform.Arch))
	}
	args = append(args, imageRef)

	stdOutWriter := newStdOutWriteCloser(eventPublisher, req.ContainerID, rootCallID)
	defer stdOutWriter.Close()

	cmd := exec.CommandContext(ctx, cr.nerdctlPath, args...)
	cmd.Stdout = stdOutWriter
	cmd.Stderr = stdOutWriter
	return cmd.Run()
}

// needsPull mirrors the docker backend: a tagged, non-"latest" image that's
// already present locally is not re-pulled. This reduces registry round-trips
// (and rate-limit exposure) and speeds up execution.
func (cr _containerRuntime) needsPull(
	ctx context.Context,
	imageRef string,
) bool {
	ref, err := reference.ParseAnyReference(strings.ToLower(imageRef))
	if err != nil {
		return true
	}
	named, err := reference.ParseNormalizedNamed(ref.String())
	if err != nil {
		return true
	}
	if tagged, ok := named.(reference.Tagged); ok && tagged.Tag() != "latest" {
		if _, err := cr.nerdctl(ctx, "image", "inspect", imageRef); err == nil {
			return false
		}
		// inspect error is expected when the image isn't present; fall through
	}
	return true
}
