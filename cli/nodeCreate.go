package main

import (
	"context"

	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	core "github.com/opctl/opctl/sdks/go/node"
)

// nodeCreate implements the node create command
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

	containerRT, err := getContainerRuntime(ctx, nodeConfig)
	if err != nil {
		return err
	}

	return newHTTPListener(
		core.New(
			ctx,
			containerRT,
			dataDir.Path(),
		),
	).
		listen(
			ctx,
			nodeConfig.ListenAddress,
		)

}
