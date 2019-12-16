package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fake.go --fake-name Fake ./ Core

import (
	"net/url"
	"os"

	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/cli/internal/clicolorer"
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/clioutput"
	"github.com/opctl/opctl/cli/internal/cliparamsatisfier"
	"github.com/opctl/opctl/cli/internal/dataresolver"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"github.com/opctl/opctl/sdks/go/node/api/client"
)

// Core exposes all cli commands
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

	apiBaseURLStr := os.Getenv("OPCTL_CLI_API_BASEURL")
	if "" == apiBaseURLStr {
		apiBaseURLStr = "http://localhost:42224/api"
	}
	apiBaseURL, err := url.Parse(apiBaseURLStr)
	if nil != err {
		panic(err)
	}

	apiClient := client.New(
		*apiBaseURL,
		&client.Opts{
			RetryLogHook: func(err error) {
				cliOutput.Attention("request resulted in a recoverable error & will be retried; error was: %v", err)
			},
		},
	)

	cliParamSatisfier := cliparamsatisfier.New(cliExiter, cliOutput)
	dataResolver := dataresolver.New(
		cliExiter,
		cliParamSatisfier,
		*apiBaseURL,
	)

	return _core{
		Eventser: newEventser(
			apiClient,
			cliExiter,
			cliOutput,
		),
		Lser: newLser(
			apiClient,
			cliExiter,
			cliOutput,
			dataResolver,
		),
		Noder: newNoder(
			cliExiter,
			nodeProvider,
		),
		Oper: newOper(
			apiClient,
			cliExiter,
			dataResolver,
		),
		Runer: newRuner(
			apiClient,
			cliColorer,
			cliExiter,
			cliOutput,
			cliParamSatisfier,
			dataResolver,
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
