// Package containerd implements opctl's ContainerRuntime by driving containerd
// through the nerdctl CLI.
//
// Unlike the docker backend (which speaks the Docker Engine API over a socket),
// this backend shells out to `nerdctl`. That choice is deliberate: nerdctl
// natively honors containerd's per-registry mirror configuration under
// `/etc/containerd/certs.d/<host>/hosts.toml` and the credential helpers
// configured in `~/.docker/config.json` (e.g. docker-credential-ecr-login).
// This lets registry pulls (Docker Hub, Quay, ...) be transparently redirected
// through a pull-through cache without any per-op image-ref changes, which the
// Docker daemon's Hub-only `registry-mirrors` setting cannot do.
//
// Registry authentication is therefore resolved by containerd/nerdctl host
// config rather than by per-op `pullCreds`; see the README for details.
package containerd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/opctl/opctl/sdks/go/node/containerruntime"
)

// containerNamePrefix is prepended to every opctl managed container so external
// tooling (and our own Delete/Kill enumeration) can recognize them. Kept
// identical to the docker backend's prefix on purpose.
const containerNamePrefix = "opctl_"

// networkName is the bridge network opctl containers are attached to so they can
// resolve one another by name.
const networkName = "opctl"

// New returns a ContainerRuntime backed by containerd via the nerdctl CLI.
//
// The nerdctl binary is resolved from $OPCTL_NERDCTL when set, otherwise from
// the first `nerdctl` found on PATH. New fails fast if the binary can't be
// located so misconfiguration surfaces at node start rather than at first run.
func New() (containerruntime.ContainerRuntime, error) {
	nerdctlPath := os.Getenv("OPCTL_NERDCTL")
	if nerdctlPath == "" {
		nerdctlPath = "nerdctl"
	}

	resolved, err := exec.LookPath(nerdctlPath)
	if err != nil {
		return nil, fmt.Errorf(
			"containerd runtime requires the nerdctl CLI; %q not found on PATH (set $OPCTL_NERDCTL to override): %w",
			nerdctlPath,
			err,
		)
	}

	return _containerRuntime{nerdctlPath: resolved}, nil
}

type _containerRuntime struct {
	nerdctlPath string
}

// nerdctl runs a nerdctl subcommand to completion, returning combined output.
func (cr _containerRuntime) nerdctl(
	ctx context.Context,
	args ...string,
) (string, error) {
	cmd := exec.CommandContext(ctx, cr.nerdctlPath, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// getContainerName converts an opctl container id into the nerdctl container
// name (prefixed so it's recognizable as opctl managed).
func getContainerName(opctlContainerID string) string {
	return fmt.Sprintf("%s%s", containerNamePrefix, opctlContainerID)
}

func (cr _containerRuntime) DeleteContainerIfExists(
	ctx context.Context,
	containerID string,
) error {
	// `rm -f` removes a running or stopped container; a missing container is not
	// an error we care about here, mirroring the docker backend's best-effort
	// cleanup semantics.
	if out, err := cr.nerdctl(ctx, "rm", "-f", getContainerName(containerID)); err != nil {
		if isNotFound(out) {
			return nil
		}
		return fmt.Errorf("unable to delete container: %w, %s", err, out)
	}
	return nil
}

// Delete removes all opctl managed containers and the opctl network.
func (cr _containerRuntime) Delete(
	ctx context.Context,
) error {
	names, err := cr.opctlContainerNames(ctx)
	if err != nil {
		return err
	}

	for _, name := range names {
		if out, err := cr.nerdctl(ctx, "rm", "-f", name); err != nil && !isNotFound(out) {
			return fmt.Errorf("unable to delete container %q: %w, %s", name, err, out)
		}
	}

	// Remove the network last; failures (not found, or still in use by a racing
	// container) are non-fatal since the network is recreated on demand.
	cr.nerdctl(ctx, "network", "rm", networkName)

	return nil
}

// Kill is equivalent to Delete for this backend (containers are stopped as part
// of removal), matching the docker backend.
func (cr _containerRuntime) Kill(
	ctx context.Context,
) error {
	return cr.Delete(ctx)
}

// opctlContainerNames lists the names of all opctl managed containers.
func (cr _containerRuntime) opctlContainerNames(
	ctx context.Context,
) ([]string, error) {
	out, err := cr.nerdctl(ctx, "ps", "--all", "--format", "{{.Names}}")
	if err != nil {
		return nil, fmt.Errorf("unable to list containers: %w, %s", err, out)
	}

	var names []string
	for _, line := range strings.Split(out, "\n") {
		name := strings.TrimSpace(line)
		if strings.HasPrefix(name, containerNamePrefix) {
			names = append(names, name)
		}
	}
	return names, nil
}

// isNotFound reports whether nerdctl output indicates a benign "no such object"
// condition (the object we asked to remove was already gone).
func isNotFound(output string) bool {
	msg := strings.ToLower(output)
	return strings.Contains(msg, "no such") || strings.Contains(msg, "not found")
}
