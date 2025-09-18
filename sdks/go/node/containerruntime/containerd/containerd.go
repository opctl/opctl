package containerd

import (
	"context"
	"fmt"
	"strings"
	"syscall"
	"time"

	containerd "github.com/containerd/containerd/v2/client"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
)

const (
	containerNamePrefix = "opctl_"
)

// New returns a ContainerRuntime backed by containerd (v2 client).
func New(
	ctx context.Context,
	address string,
) (containerruntime.ContainerRuntime, error) {
	fmt.Println("containerd: initializing...")

	cli, err := containerd.New(
		address,
		containerd.WithDefaultNamespace("default"),
	)
	if err != nil {
		return nil, fmt.Errorf("containerd: connect failed: %w", err)
	}
	fmt.Println("containerd: connecting...")

	// Add timeout to prevent hanging
	serveCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	isConnected, err := cli.IsServing(serveCtx)
	if err != nil {
		return nil, fmt.Errorf("error connecting to containerd: %w", err)
	}
	fmt.Printf("containerd: connected: %t\n", isConnected)

	if !isConnected {
		return nil, fmt.Errorf("containerd not connected")
	}

	rc, err := newRunContainer(cli)
	if err != nil {
		return nil, err
	}

	return _containerRuntime{
		client:       cli,
		runContainer: rc,
	}, nil
}

type _containerRuntime struct {
	client *containerd.Client
	runContainer
}

func (cr _containerRuntime) Delete(
	ctx context.Context,
) error {
	containers, err := cr.client.Containers(ctx)
	if err != nil {
		return fmt.Errorf("containerd: list containers: %w", err)
	}

	var retErr error
	for _, c := range containers {
		name := c.ID()
		if !strings.HasPrefix(name, containerNamePrefix) {
			continue
		}

		// try to stop task if running
		if task, err := c.Task(ctx, nil); err == nil {
			_ = task.Kill(ctx, syscall.SIGTERM)
			_ = task.Kill(ctx, syscall.SIGKILL)
			_, _ = task.Delete(ctx, containerd.WithProcessKill)
		}
		// remove container and its snapshot
		if err := c.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
			retErr = fmt.Errorf("containerd: delete container %s: %w", name, err)
		}
	}

	return retErr
}

func (cr _containerRuntime) DeleteContainerIfExists(
	ctx context.Context,
	containerID string,
) error {
	name := getContainerName(containerID)

	c, err := cr.client.LoadContainer(ctx, name)
	if err != nil {
		// treat not found as success
		return nil
	}

	if task, err := c.Task(ctx, nil); err == nil {
		_ = task.Kill(ctx, syscall.SIGTERM)
		_ = task.Kill(ctx, syscall.SIGKILL)
		_, _ = task.Delete(ctx, containerd.WithProcessKill)
	}

	if err := c.Delete(ctx, containerd.WithSnapshotCleanup); err != nil {
		return fmt.Errorf("containerd: delete container %s: %w", name, err)
	}

	return nil
}

func (cr _containerRuntime) Kill(
	ctx context.Context,
) error {
	return cr.Delete(ctx)
}

func getContainerName(opctlContainerID string) string {
	return fmt.Sprintf("%s%s", containerNamePrefix, opctlContainerID)
}
