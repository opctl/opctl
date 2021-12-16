package main

import (
	"context"

	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
)

// nodeKill implements the node kill command
func nodeKill(
	ctx context.Context,
	nodeConfig local.NodeConfig,
) error {
	containerRT, err := getContainerRuntime(ctx, nodeConfig)
	if err != nil {
		return err
	}

	if err := local.New(nodeConfig).KillNodeIfExists(""); err != nil {
		return err
	}

	return containerRT.Kill(ctx)

}
