package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/lazylocalnode"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
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
	nodeProviderOpts local.NodeCreateOpts,
) Core {
	cliParamSatisfier := cliparamsatisfier.New(cliOutput)

	nodeProvider := local.New(nodeProviderOpts)

	opNode := lazylocalnode.New(nodeProvider)

	dataResolver := dataresolver.New(
		cliParamSatisfier,
		opNode,
	)

	return _core{
		Auther: newAuther(
			dataResolver,
			opNode,
		),
		Eventser: newEventser(
			cliOutput,
			opNode,
		),
		Lser: newLser(
			cliOutput,
			dataResolver,
		),
		Noder: newNoder(nodeProvider),
		Oper: newOper(
			dataResolver,
			opNode,
		),
		Runer: newRuner(
			cliOutput,
			cliParamSatisfier,
			dataResolver,
			opNode,
		),
		SelfUpdater: newSelfUpdater(nodeProvider),
		UIer: newUIer(
			dataResolver,
			opNode,
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
