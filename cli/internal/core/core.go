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
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// Core exposes all cli commands
//counterfeiter:generate -o fakes/core.go . Core
type Core interface {
	Auther
	Eventser
	Lser
	Noder
	Oper
	Runer
	SelfUpdater
	UIer
	cliexiter.CliExiter
}

// New returns initialized cli core
func New(
	cliColorer clicolorer.CliColorer,
	nodeProvider nodeprovider.NodeProvider,
) Core {
	_os := ios.New()
	cliOutput := clioutput.New(cliColorer, os.Stderr, os.Stdout)
	cliExiter := cliexiter.New(cliOutput, _os)
	cliParamSatisfier := cliparamsatisfier.New(cliOutput)

	dataResolver := dataresolver.New(
		cliParamSatisfier,
		nodeProvider,
	)

	return _core{
		Auther: newAuther(
			dataResolver,
			nodeProvider,
		),
		Eventser: newEventser(
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
		UIer: newUIer(
			dataResolver,
			nodeProvider,
		),
		CliExiter: cliExiter,
	}
}

type _core struct {
	Auther
	Eventser
	Lser
	Noder
	Oper
	Runer
	SelfUpdater
	UIer
	cliexiter.CliExiter
}
