package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/skratchdot/open-golang/open"
)

func openUI(
	nodeProvider nodeprovider.NodeProvider,
	cliParamSatisfier cliparamsatisfier.CLIParamSatisfier,
	mountRef string,
) error {
	var resolvedMount string
	var err error
	if strings.HasPrefix(mountRef, ".") {
		// treat dot paths as regular rel paths
		resolvedMount, err = filepath.Abs(mountRef)
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
			mountRef,
			nil,
		)
		if nil != err {
			return err
		}

		resolvedMount = mountHandle.Ref()
	}

	open.Run(
		fmt.Sprintf("http://localhost:42224?mount=%s", url.QueryEscape(resolvedMount)),
	)
	return nil
}
