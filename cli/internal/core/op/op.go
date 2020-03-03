package op

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// Op exposes the "op" sub command
//counterfeiter:generate -o fakes/op.go . Op
type Op interface {
	Creater
	Installer
	Killer
	Validater
}

// New returns an initialized "op" sub command
func New(
	cliExiter cliexiter.CliExiter,
	dataResolver dataresolver.DataResolver,
	nodeProvider nodeprovider.NodeProvider,
) Op {
	return _op{
		Creater: newCreater(
			cliExiter,
		),
		Installer: newInstaller(
			cliExiter,
			dataResolver,
		),
		Killer: newKiller(
			cliExiter,
			nodeProvider,
		),
		Validater: newValidater(
			cliExiter,
			dataResolver,
		),
	}
}

type _op struct {
	Creater
	Installer
	Killer
	Validater
}
