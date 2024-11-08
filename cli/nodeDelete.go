package main

import (
	"context"
	"os"

	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/node/dns"
)

// nodeDelete implements the node delete command
func nodeDelete(
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

	err = containerRT.Delete(ctx)
	if err != nil {
		return err
	}

	if err := np.KillNodeIfExists(); err != nil {
		return err
	}

	if err := dns.Delete(); err != nil {
		return err
	}

	return os.RemoveAll(nodeConfig.DataDir)

}
