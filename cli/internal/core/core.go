package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
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
}

// New returns initialized cli core
func New(
	cliOutput clioutput.CliOutput,
	nodeProvider nodeprovider.NodeProvider,
) Core {
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
			cliOutput,
			dataResolver,
		),
		Noder: newNoder(nodeProvider),
		Oper: newOper(
			dataResolver,
			nodeProvider,
		),
		Runer: newRuner(
			cliOutput,
			cliParamSatisfier,
			dataResolver,
			nodeProvider,
		),
		SelfUpdater: newSelfUpdater(nodeProvider),
		UIer: newUIer(
			dataResolver,
			nodeProvider,
		),
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
}
