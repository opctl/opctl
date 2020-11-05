package op

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
)

// Installer exposes the "op install" sub command
type Installer interface {
	Install(
		ctx context.Context,
		path,
		opRef,
		username,
		password string,
	)
}

// newInstaller returns an initialized "op install" sub command
func newInstaller(
	cliExiter cliexiter.CliExiter,
	dataResolver dataresolver.DataResolver,
) Installer {
	return _installer{
		cliExiter:    cliExiter,
		dataResolver: dataResolver,
	}
}

type _installer struct {
	cliExiter    cliexiter.CliExiter
	dataResolver dataresolver.DataResolver
}

func (ivkr _installer) Install(
	ctx context.Context,
	path,
	opRef,
	username,
	password string,
) {
	// install the whole pkg in case relative (intra pkg) refs exist
	opRefParts := strings.SplitN(opRef, "#", 2)
	var pkgRef string
	if len(opRefParts) == 1 {
		pkgRef = opRefParts[0]
	} else {
		if verAndPathParts := strings.SplitN(opRefParts[1], "/", 2); len(verAndPathParts) != 1 {
			pkgRef = fmt.Sprintf("%s#%s", opRefParts[0], verAndPathParts[0])
		}
	}

	opDirHandle := ivkr.dataResolver.Resolve(
		pkgRef,
		&model.Creds{
			Username: username,
			Password: password,
		},
	)

	if err := opspec.Install(
		ctx,
		filepath.Join(path, pkgRef),
		opDirHandle,
	); nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
	}

}
