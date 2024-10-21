package main

import (
	"context"
	"fmt"
	"os"

	"syscall"

	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/node/dns"
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

	np, err := local.New(nodeConfig)
	if err != nil {
		return err
	}

	if os.Geteuid() != 0 {
		err := syscall.Seteuid(0)
		if err != nil {
			return fmt.Errorf("re-run command with sudo: %s", err.Error())
		}
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
