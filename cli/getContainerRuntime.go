package main

import (
	"context"

	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/k8s"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/qemu"
)

func getContainerRuntime(
	ctx context.Context,
	nodeConfig local.NodeConfig,
) (containerruntime.ContainerRuntime, error) {
	if nodeConfig.ContainerRuntime == "k8s" {
		return k8s.New()
	} else if nodeConfig.ContainerRuntime == "qemu" {
		return qemu.New(ctx, false)
	} else {
		return docker.New(ctx, "unix:///var/run/docker.sock")
	}
}
