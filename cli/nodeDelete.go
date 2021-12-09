package main

import (
	"context"
	"os"

	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/k8s"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/qemu"
)

// node command
func nodeDelete(
	ctx context.Context,
	nodeConfig local.NodeConfig,
) error {

	var containerRuntime containerruntime.ContainerRuntime
	var err error
	if nodeConfig.ContainerRuntime == "k8s" {
		containerRuntime, err = k8s.New()
		if err != nil {
			return err
		}
	} else if nodeConfig.ContainerRuntime == "qemu" {
		containerRuntime, err = qemu.New(ctx, false)
		if err != nil {
			return err
		}
	} else {
		containerRuntime, err = docker.New(ctx, "unix:///var/run/docker.sock")
		if err != nil {
			return err
		}
	}

	if err := local.New(nodeConfig).KillNodeIfExists(""); err != nil {
		return err
	}

	if err := os.RemoveAll(nodeConfig.DataDir); err != nil {
		return err
	}

	return containerRuntime.Delete(ctx)

}
