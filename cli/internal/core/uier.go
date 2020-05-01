package core

import (
	"fmt"
	"path/filepath"

	"net/url"

	"strings"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/skratchdot/open-golang/open"
)

// Uier exposes the "ui" command
type UIer interface {
	UI(
		mountRef string,
	)
}

// newUIer returns an initialized "ui" command
func newUIer(
	cliExiter cliexiter.CliExiter,
	dataResolver dataresolver.DataResolver,
	nodeProvider nodeprovider.NodeProvider,
) UIer {
	return _uier{
		cliExiter:    cliExiter,
		dataResolver: dataResolver,
		nodeProvider: nodeProvider,
	}
}

type _uier struct {
	cliExiter    cliexiter.CliExiter
	dataResolver dataresolver.DataResolver
	nodeProvider nodeprovider.NodeProvider
}

func (ivkr _uier) UI(
	mountRef string,
) {

	_, err := ivkr.nodeProvider.CreateNodeIfNotExists()
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	var resolvedMount string
	if strings.HasPrefix(mountRef, ".") {
		// treat dot paths as regular rel paths
		resolvedMount, err = filepath.Abs(mountRef)
		if nil != err {
			ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
			return // support fake exiter
		}
	} else {
		// otherwise use same resolution as run
		mountHandle := ivkr.dataResolver.Resolve(
			mountRef,
			nil,
		)
		resolvedMount = mountHandle.Ref()
	}

	webUIURL := fmt.Sprintf("http://localhost:42224?mount=%s", url.QueryEscape(resolvedMount))

	err = open.Run(webUIURL)
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	ivkr.cliExiter.Exit(cliexiter.ExitReq{
		Message: fmt.Sprint("Opctl web UI opened!\n"),
		Code:    0,
	})

}
