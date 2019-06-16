package core

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Core

import (
	"context"
	"io"
	"net/url"
	"os"

	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/nodeprovider"
	"github.com/opctl/opctl/nodeprovider/local"
	"github.com/opctl/opctl/sdk/go/node/api/client"
	op "github.com/opctl/opctl/sdk/go/opspec"
	dotyml "github.com/opctl/opctl/sdk/go/opspec/opfile"
	"github.com/opctl/opctl/util/clicolorer"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/clioutput"
	"github.com/opctl/opctl/util/cliparamsatisfier"
	"github.com/opctl/opctl/util/updater"
)

type Core interface {
	Events(
		ctx context.Context,
	)

	NodeCreate(
		opts NodeCreateOpts,
	)

	NodeKill()

	Ls(
		ctx context.Context,
		dirRef string,
	)

	OpCreate(
		path string,
		description string,
		name string,
	)

	OpInstall(
		ctx context.Context,
		path,
		opRef,
		username,
		password string,
	)

	OpKill(
		ctx context.Context,
		opId string,
	)

	OpValidate(
		ctx context.Context,
		opRef string,
	)

	Run(
		ctx context.Context,
		opRef string,
		opts *RunOpts,
	)

	SelfUpdate(
		channel string,
	)
}

func New(
	cliColorer clicolorer.CliColorer,
) Core {

	_os := ios.New()

	cliOutput := clioutput.New(cliColorer, os.Stderr, os.Stdout)
	cliExiter := cliexiter.New(cliOutput, _os)
	cliParamSatisfier := cliparamsatisfier.New(cliExiter, cliOutput)

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

	return &_core{
		cliColorer:              cliColorer,
		cliExiter:               cliExiter,
		cliOutput:               cliOutput,
		cliParamSatisfier:       cliParamSatisfier,
		opDotYmlGetter:          dotyml.NewGetter(),
		ioutil:                  iioutil.New(),
		nodeProvider:            local.New(),
		nodeReachabilityEnsurer: newNodeReachabilityEnsurer(cliExiter),
		opCreator:               op.NewCreator(),
		opInstaller:             op.NewInstaller(),
		opLister:                op.NewLister(),
		apiClient:               apiClient,
		opValidator:             op.NewValidator(),
		os:                      _os,
		dataResolver:            newDataResolver(cliExiter, cliParamSatisfier, *apiBaseURL),
		updater:                 updater.New(),
		writer:                  os.Stdout,
	}

}

type _core struct {
	cliColorer              clicolorer.CliColorer
	cliExiter               cliexiter.CliExiter
	cliOutput               clioutput.CliOutput
	cliParamSatisfier       cliparamsatisfier.CLIParamSatisfier
	opDotYmlGetter          dotyml.Getter
	ioutil                  iioutil.IIOUtil
	nodeProvider            nodeprovider.NodeProvider
	nodeReachabilityEnsurer nodeReachabilityEnsurer
	opCreator               op.Creator
	opInstaller             op.Installer
	opLister                op.Lister
	apiClient               client.Client
	opValidator             op.Validator
	os                      ios.IOS
	dataResolver            dataResolver
	updater                 updater.Updater
	writer                  io.Writer
}
