package op

import (
	"context"

	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/model"
	opspec "github.com/opctl/opctl/sdks/go/opspec"
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
		opInstaller:  opspec.NewInstaller(),
	}
}

type _installer struct {
	cliExiter    cliexiter.CliExiter
	dataResolver dataresolver.DataResolver
	opInstaller  opspec.Installer
}

func (ivkr _installer) Install(
	ctx context.Context,
	path,
	opRef,
	username,
	password string,
) {

	opDirHandle := ivkr.dataResolver.Resolve(
		opRef,
		&model.PullCreds{
			Username: username,
			Password: password,
		},
	)

	if err := ivkr.opInstaller.Install(
		ctx,
		path,
		opDirHandle,
	); nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
	}

}
