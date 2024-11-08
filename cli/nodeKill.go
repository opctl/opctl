package main

import (
	"context"

	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/node/dns"
)

// nodeKill implements the node kill command
func nodeKill(
	ctx context.Context,
	nodeConfig local.NodeConfig,
) error {
	if err := ensureEuid0(); err != nil {
		return err
	}

	containerRT, err := getContainerRuntime(ctx, nodeConfig)
	if err != nil {
		return err
	}

	np, err := local.New(nodeConfig)
	if err != nil {
		return err
	}

	if err := np.KillNodeIfExists(); err != nil {
		return err
	}

	if err := dns.Delete(); err != nil {
		return err
	}

	return containerRT.Kill(ctx)

}
