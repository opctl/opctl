package main

import (
	"context"

	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	core "github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/k8s"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/qemu"
)

// node command
func nodeCreate(
	ctx context.Context,
	nodeConfig local.NodeConfig,
) error {
	dataDir, err := datadir.New(nodeConfig.DataDir)
	if err != nil {
		return err
	}

	if err := dataDir.InitAndLock(); err != nil {
		return err
	}

	var containerRuntime containerruntime.ContainerRuntime
	if nodeConfig.ContainerRuntime == "k8s" {
		containerRuntime, err = k8s.New()
		if err != nil {
			return err
		}
	} else if nodeConfig.ContainerRuntime == "qemu" {
		containerRuntime, err = qemu.New(ctx, true)
		if err != nil {
			return err
		}
	} else {
		containerRuntime, err = docker.New(ctx, "unix:///var/run/docker.sock")
		if err != nil {
			return err
		}
	}

	return newHTTPListener(
		core.New(
			ctx,
			containerRuntime,
			dataDir.Path(),
		),
	).
		listen(
			ctx,
			nodeConfig.ListenAddress,
		)

}
