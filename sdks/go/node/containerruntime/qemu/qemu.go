package qemu

import (
	"context"
	"fmt"
	"io"
	"runtime"

	"github.com/abiosoft/colima/environment"
	"github.com/abiosoft/colima/environment/host"
	"github.com/abiosoft/colima/environment/vm/lima"

	"github.com/abiosoft/colima/config"
	colimadocker "github.com/abiosoft/colima/environment/container/docker"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
	"github.com/pbnjay/memory"
)

func New(
	ctx context.Context,
	waitUntilReady bool,
) (containerruntime.ContainerRuntime, error) {
	cr := _containerRuntime{
		vm: lima.New(host.New()),
	}

	if err := host.IsInstalled(cr.vm); err != nil {
		// ensure lima present
		return nil, fmt.Errorf("dependency check failed for VM: %w", err)
	}

	if waitUntilReady {
		if _, err := cr.getDockerContainerRuntime(ctx); err != nil {
			return nil, err
		}
	}

	return cr, nil
}

type _containerRuntime struct {
	vm environment.VM
}

func (cr _containerRuntime) Delete(
	ctx context.Context,
) error {
	return cr.vm.Teardown()
}

func (cr _containerRuntime) DeleteContainerIfExists(
	ctx context.Context,
	containerID string,
) error {
	if !cr.vm.Running() {
		return nil
	}

	dockerCR, err := cr.getDockerContainerRuntime(ctx)
	if err != nil {
		return err
	}

	return dockerCR.DeleteContainerIfExists(ctx, containerID)
}

func (cr _containerRuntime) Kill(
	ctx context.Context,
) error {
	if !cr.vm.Running() {
		return nil
	}

	dockerCR, err := cr.getDockerContainerRuntime(ctx)
	if err != nil {
		return err
	}

	if err := dockerCR.Kill(ctx); err != nil {
		return err
	}

	return cr.vm.Stop()
}

// RunContainer creates, starts, and waits on a container. ExitCode &/Or an error will be returned
func (cr _containerRuntime) RunContainer(
	ctx context.Context,
	req *model.ContainerCall,
	// @TODO: get rid of in combination with eventPublisher
	rootCallID string,
	// @TODO: get rid of this; just use stdout/stderr
	eventPublisher pubsub.EventPublisher,
	stdout io.WriteCloser,
	stderr io.WriteCloser,
) (*int64, error) {
	dockerCR, err := cr.getDockerContainerRuntime(ctx)
	if err != nil {
		return nil, err
	}

	return dockerCR.RunContainer(ctx, req, rootCallID, eventPublisher, stdout, stderr)
}

func (cr _containerRuntime) getDockerContainerRuntime(
	ctx context.Context,
) (containerruntime.ContainerRuntime, error) {
	if !cr.vm.Running() {
		err := cr.vm.Start(config.Config{
			Runtime: "docker",
			VM: config.VM{
				Arch: string(environment.Arch(runtime.GOARCH).Value()),
				// allocate 2/3 available CPU's
				CPU:  runtime.NumCPU() * 2 / 3,
				Disk: 60,
				// allocate 2/3 available memory
				Memory: int(memory.TotalMemory()) * 2 / 3e9,
			},
		})
		if err != nil {
			return nil, err
		}

		// start docker, sleep while it creates the socket, and grant access to socket
		if err := cr.vm.Run("sudo", "sh", "-ce", "service docker start && sleep 2 && chmod 0666 /var/run/docker.sock"); err != nil {
			return nil, fmt.Errorf("error adding VM user to docker group: %w", err)
		}
	}
	return docker.New(ctx, fmt.Sprintf("unix://%s", colimadocker.HostSocketFile()))
}
