package node

import (
	"context"

	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/docker"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/embedded"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/k8s"
)

func getContainerRuntime(
	ctx context.Context,
	nodeConfig local.NodeConfig,
) (containerruntime.ContainerRuntime, error) {
	if nodeConfig.ContainerRuntime == "k8s" {
		return k8s.New()
	} else if nodeConfig.ContainerRuntime == "embedded" {
		return embedded.New(ctx, nodeConfig.DataDir)
	} else {
		return docker.New(ctx, "unix:///var/run/docker.sock")
	}
}
