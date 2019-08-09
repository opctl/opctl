package op

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fake.go --fake-name Fake ./ Op

import (
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

// Op exposes the "op" sub command
type Op interface {
	Creater
	Installer
	Killer
	Validater
}

// New returns an initialized "op" sub command
func New(
	apiClient client.Client,
	cliExiter cliexiter.CliExiter,
	dataResolver dataresolver.DataResolver,
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
			apiClient,
			cliExiter,
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
