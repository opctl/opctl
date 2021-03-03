package main

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/skratchdot/open-golang/open"
)

// ui implements "ui" command
func ui(
	ctx context.Context,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeProvider nodeprovider.NodeProvider,
	listenAddress string,
	mountRefArg string,
) error {
	var resolvedMount string
	node, err := nodeProvider.CreateNodeIfNotExists(ctx)
	if err != nil {
		return err
	}

	if strings.HasPrefix(mountRefArg, ".") {
		// treat dot paths as regular rel paths
		resolvedMount, err = filepath.Abs(mountRefArg)
		if nil != err {
			return err
		}
	} else {
		dataResolver := dataresolver.New(
			cliParamSatisfier,
			node,
		)

		// otherwise use same resolution as run
		mountHandle, err := dataResolver.Resolve(
			ctx,
			mountRefArg,
			nil,
		)
		if nil != err {
			return err
		}

		resolvedMount = mountHandle.Ref()
	}

	return open.Run(
		fmt.Sprintf("http://%s/?mount=%s", listenAddress, url.QueryEscape(resolvedMount)),
	)
}
