package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"os"

	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
)

// Core exposes all cli commands
//counterfeiter:generate -o fakes/core.go . Core
type Core interface {
	Eventser
	Lser
	Noder
	Oper
	Runer
	SelfUpdater
}

// New returns initialized cli core
func New(
	cliColorer clicolorer.CliColorer,
) Core {
	_os := ios.New()
	cliOutput := clioutput.New(cliColorer, os.Stderr, os.Stdout)
	cliExiter := cliexiter.New(cliOutput, _os)
	nodeProvider := local.New()

	cliParamSatisfier := cliparamsatisfier.New(cliExiter, cliOutput)
	dataResolver := dataresolver.New(
		cliExiter,
		cliParamSatisfier,
		nodeProvider,
	)

	return _core{
		Eventser: newEventser(
			cliExiter,
			cliOutput,
			nodeProvider,
		),
		Lser: newLser(
			cliExiter,
			cliOutput,
			dataResolver,
		),
		Noder: newNoder(
			cliExiter,
			nodeProvider,
		),
		Oper: newOper(
			cliExiter,
			dataResolver,
			nodeProvider,
		),
		Runer: newRuner(
			cliColorer,
			cliExiter,
			cliOutput,
			cliParamSatisfier,
			dataResolver,
			nodeProvider,
		),
		SelfUpdater: newSelfUpdater(
			cliExiter,
			nodeProvider,
		),
	}
}

type _core struct {
	Eventser
	Lser
	Noder
	Oper
	Runer
	SelfUpdater
}
