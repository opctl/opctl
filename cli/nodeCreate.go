package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/euid0"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/cli/internal/pidfile"
	core "github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/api"
	"github.com/opctl/opctl/sdks/go/node/datadir"
	"golang.org/x/sync/errgroup"
)

// nodeCreate implements the node create command
func nodeCreate(
	ctx context.Context,
	nodeConfig local.NodeConfig,
) error {
	if err := euid0.Ensure(); err != nil {
		return err
	}

	cliColorer := clicolorer.New()

	dataDir, err := datadir.New(nodeConfig.DataDir)
	if err != nil {
		return err
	}

	gotLock, err := pidfile.TryGetLock(
		ctx,
		nodeConfig.DataDir,
	)
	if err != nil {
		return err
	}

	if !gotLock {
		return fmt.Errorf("node already running; to kill use \"sudo opctl node kill\"")
	}

	containerRT, err := getContainerRuntime(ctx, nodeConfig)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)

	// catch signals to ensure shutdown properly happens
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	eg.Go(
		func() error {
			fmt.Println(
				cliColorer.Info(fmt.Sprintf("opctl API listening at %s", nodeConfig.APIListenAddress)),
			)

			c, err := core.New(
				ctx,
				containerRT,
				dataDir.Path(),
			)
			if err != nil {
				return err
			}

			return api.Listen(
				ctx,
				nodeConfig.APIListenAddress,
				c,
			)
		},
	)

	err = eg.Wait()
	stop()

	return err

}
