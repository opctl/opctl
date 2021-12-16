package main

import (
	"context"
	"os"

	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
)

// nodeDelete implements the node delete command
func nodeDelete(
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

	if err := os.RemoveAll(nodeConfig.DataDir); err != nil {
		return err
	}

	return containerRT.Delete(ctx)

}
