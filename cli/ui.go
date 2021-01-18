package main

import (
	"fmt"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/skratchdot/open-golang/open"
	"net/url"
	"path/filepath"
	"strings"
)

// ui implements "ui" command
func ui(
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	nodeProvider nodeprovider.NodeProvider,
	mountRefArg string,
) error {
	var resolvedMount string
	var err error
	if strings.HasPrefix(mountRefArg, ".") {
		// treat dot paths as regular rel paths
		resolvedMount, err = filepath.Abs(mountRefArg)
		if nil != err {
			return err
		}
	} else {
		node, err := nodeProvider.CreateNodeIfNotExists()
		if err != nil {
			return err
		}

		dataResolver := dataresolver.New(
			cliParamSatisfier,
			node,
		)

		// otherwise use same resolution as run
		mountHandle, err := dataResolver.Resolve(
			mountRefArg,
			nil,
		)
		if nil != err {
			return err
		}

		resolvedMount = mountHandle.Ref()
	}

	return open.Run(
		fmt.Sprintf("http://localhost:42224?mount=%s", url.QueryEscape(resolvedMount)),
	)
}
