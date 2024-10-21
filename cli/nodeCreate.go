package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	core "github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/api"
	"github.com/opctl/opctl/sdks/go/node/datadir"
	"github.com/opctl/opctl/sdks/go/node/dns"
	"golang.org/x/sync/errgroup"
)

// nodeCreate implements the node create command
func nodeCreate(
	ctx context.Context,
	nodeConfig local.NodeConfig,
) error {
	if os.Geteuid() != 0 {
		return errors.New("re-run command with sudo")
	}

	cliColorer := clicolorer.New()

	dataDir, err := datadir.New(nodeConfig.DataDir)
	if err != nil {
		return err
	}

	unlock, err := dataDir.InitAndLock()
	if err != nil {
		return err
	}

	defer unlock()

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

			return api.Listen(
				ctx,
				nodeConfig.APIListenAddress,
				core.New(
					ctx,
					containerRT,
					dataDir.Path(),
				),
			)
		},
	)

	eg.Go(
		func() error {
			fmt.Println(
				cliColorer.Info(fmt.Sprintf("opctl DNS listening at %s", nodeConfig.DNSListenAddress)),
			)

			return dns.Listen(
				ctx,
				nodeConfig.DNSListenAddress,
			)
		},
	)

	err = eg.Wait()
	stop()

	return err

}
