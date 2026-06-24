package containerd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/dns"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
	"golang.org/x/sync/errgroup"
)

// RunContainer creates, starts, and waits on a container, returning its exit
// code and/or an error. It mirrors the docker backend's lifecycle (ensure
// network, pull/load image, create, start, stream stdout/stderr, register DNS
// for named containers, wait) but drives containerd via nerdctl.
func (cr _containerRuntime) RunContainer(
	ctx context.Context,
	req *model.ContainerCall,
	rootCallID string,
	eventPublisher pubsub.EventPublisher,
	stdout io.WriteCloser,
	stderr io.WriteCloser,
) (*int64, error) {
	defer stdout.Close()
	defer stderr.Close()

	// ensure the opctl network exists so named containers resolve one another
	if err := cr.ensureNetworkExists(ctx); err != nil {
		return nil, err
	}

	name := getContainerName(req.ContainerID)

	// remove any stale container left behind by a prior run sharing this id
	cr.nerdctl(context.Background(), "rm", "-f", name)
	defer func() {
		// always clean up, even after cancellation, using a fresh context
		cr.nerdctl(context.Background(), "rm", "-f", name)
	}()

	var imageErr error
	if req.Image.Src != nil {
		imageErr = cr.loadImage(ctx, req)
	} else {
		// don't err yet; the image might be cached, which we allow for offline use
		imageErr = cr.pullImage(ctx, req, rootCallID, eventPublisher)
	}

	createArgs, err := constructCreateArgs(req, name)
	if err != nil {
		return nil, errors.Join(imageErr, err)
	}

	if out, createErr := cr.nerdctl(ctx, createArgs...); createErr != nil {
		select {
		case <-ctx.Done():
			// we got killed
			return nil, nil
		default:
			return nil, errors.Join(imageErr, fmt.Errorf("unable to create container: %w, %s", createErr, out))
		}
	}

	if out, err := cr.nerdctl(ctx, "start", name); err != nil {
		return nil, fmt.Errorf("unable to start container: %w, %s", err, out)
	}

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return cr.streamLogs(egCtx, name, stdout, stderr)
	})

	// register named containers with opctl's DNS so they're resolvable by name
	if req.Name != nil {
		ip, err := cr.containerIP(ctx, name)
		if err != nil {
			return nil, err
		}
		if ip != "" {
			defer dns.UnregisterName(*req.Name, ip)
			if err := dns.RegisterName(ctx, *req.Name, ip); err != nil {
				return nil, err
			}
		}
	}

	exitCode, waitErr := cr.waitContainer(ctx, name)

	// ensure all stdout/stderr is read before returning
	logErr := eg.Wait()

	return exitCode, errors.Join(waitErr, logErr)
}

// ensureNetworkExists creates the opctl bridge network, tolerating the common
// case where it already exists (and races between concurrent container starts).
func (cr _containerRuntime) ensureNetworkExists(
	ctx context.Context,
) error {
	out, err := cr.nerdctl(ctx, "network", "create", networkName)
	if err != nil && !strings.Contains(strings.ToLower(out), "already exists") {
		return fmt.Errorf("unable to create network: %w, %s", err, out)
	}
	return nil
}

// constructCreateArgs builds the `nerdctl create` argument vector for a
// container call. Ordering of env vars, mounts, and ports is sorted so the
// result is deterministic (and testable).
func constructCreateArgs(
	req *model.ContainerCall,
	name string,
) ([]string, error) {
	args := []string{
		"create",
		"--name", name,
		"--network", networkName,
		// support docker-in-docker, matching the docker backend's HostConfig
		"--privileged",
	}

	if req.Name != nil {
		// network alias makes the container resolvable by its op-declared name
		args = append(args, "--network-alias", *req.Name)
	}

	if req.WorkDir != "" {
		args = append(args, "--workdir", req.WorkDir)
	}

	// env vars: op-declared values plus propagated proxy vars (proxy vars never
	// clobber op-declared ones; see containerruntime.ProxyEnvVars).
	env := map[string]string{}
	for k, v := range req.EnvVars {
		env[k] = v
	}
	for k, v := range containerruntime.ProxyEnvVars(req.EnvVars) {
		env[k] = v
	}
	for _, k := range sortedKeys(env) {
		args = append(args, "--env", fmt.Sprintf("%s=%s", k, env[k]))
	}

	// bind mounts: files, then dirs, then sockets (matching the docker backend),
	// each group sorted by container path for determinism.
	for _, containerPath := range sortedKeys(req.Files) {
		args = append(args, "--volume", fmt.Sprintf("%s:%s", req.Files[containerPath], containerPath))
	}
	for _, containerPath := range sortedKeys(req.Dirs) {
		args = append(args, "--volume", fmt.Sprintf("%s:%s", req.Dirs[containerPath], containerPath))
	}
	for _, containerSocket := range sortedKeys(req.Sockets) {
		hostSocket := req.Sockets[containerSocket]
		// only unix sockets (identified naively by a path separator) are bind
		// mounted, matching the docker backend.
		if strings.ContainsAny(hostSocket, `/\`) {
			args = append(args, "--volume", fmt.Sprintf("%s:%s", hostSocket, containerSocket))
		}
	}

	// port bindings: hostPort:containerPort
	for _, containerPort := range sortedKeys(req.Ports) {
		args = append(args, "--publish", fmt.Sprintf("%s:%s", req.Ports[containerPort], containerPort))
	}

	if req.Image == nil || req.Image.Ref == nil {
		return nil, errors.New("image ref is required")
	}
	args = append(args, *req.Image.Ref)
	args = append(args, req.Cmd...)

	return args, nil
}

// streamLogs follows the container's logs, writing stdout and stderr to the
// supplied writers. nerdctl demultiplexes the two streams (the container is not
// allocated a TTY), so they stay separated. Errors are treated as benign: the
// container may already be gone, or the context cancelled.
func (cr _containerRuntime) streamLogs(
	ctx context.Context,
	name string,
	stdout io.Writer,
	stderr io.Writer,
) error {
	cmd := exec.CommandContext(ctx, cr.nerdctlPath, "logs", "--follow", name)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Run()
	return nil
}

// containerIP returns the container's IP on the opctl network.
func (cr _containerRuntime) containerIP(
	ctx context.Context,
	name string,
) (string, error) {
	out, err := cr.nerdctl(
		ctx,
		"inspect",
		"--format", "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}",
		name,
	)
	if err != nil {
		return "", fmt.Errorf("unable to inspect container: %w, %s", err, out)
	}
	return strings.TrimSpace(out), nil
}

// waitContainer blocks until the container exits and returns its exit code.
func (cr _containerRuntime) waitContainer(
	ctx context.Context,
	name string,
) (*int64, error) {
	out, err := cr.nerdctl(ctx, "wait", name)
	if err != nil {
		if ctx.Err() != nil {
			// killed; not a real wait failure
			return nil, nil
		}
		return nil, fmt.Errorf("error waiting on container: %w, %s", err, out)
	}

	exitCode, parseErr := strconv.ParseInt(lastNonEmptyLine(out), 10, 64)
	if parseErr != nil {
		return nil, fmt.Errorf("unable to parse container exit code from %q: %w", out, parseErr)
	}
	return &exitCode, nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func lastNonEmptyLine(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if line := strings.TrimSpace(lines[i]); line != "" {
			return line
		}
	}
	return ""
}
