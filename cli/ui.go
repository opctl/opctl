package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/cli/internal/oppath"
	"github.com/skratchdot/open-golang/open"
)

// ui implements "ui" command
func ui(
	ctx context.Context,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeConfig local.NodeConfig,
	listenAddress string,
	mountRefArg string,
) error {
	var resolvedMount string

	np, err := local.New(nodeConfig)
	if err != nil {
		return err
	}

	node, err := np.CreateNodeIfNotExists(ctx)
	if err != nil {
		return err
	}

	if strings.HasPrefix(mountRefArg, ".") {
		// treat dot paths as regular rel paths
		resolvedMount, err = filepath.Abs(mountRefArg)
		if err != nil {
			return err
		}
	} else {
		dataResolver := dataresolver.New(
			cliParamSatisfier,
			node,
		)

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		opPath, err := oppath.Get(
			ctx,
			cwd,
			dataResolver,
			node,
		)
		if err != nil {
			return err
		}

		// otherwise use same resolution as run
		mountHandle, err := dataResolver.Resolve(
			ctx,
			mountRefArg,
			opPath,
			nil,
		)
		if err != nil {
			return err
		}

		resolvedMount = mountHandle.Ref()
	}

	return open.Run(
		fmt.Sprintf("http://%s/?mount=%s", listenAddress, url.QueryEscape(resolvedMount)),
	)
}
