package core

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Core

import (
	"context"
	"github.com/golang-interfaces/iioutil"
	"github.com/golang-interfaces/ios"
	"github.com/opctl/opctl/nodeprovider"
	"github.com/opctl/opctl/nodeprovider/local"
	"github.com/opctl/opctl/util/clicolorer"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/clioutput"
	"github.com/opctl/opctl/util/cliparamsatisfier"
	"github.com/opctl/opctl/util/updater"
	"github.com/opspec-io/sdk-golang/node/api/client"
	"github.com/opspec-io/sdk-golang/pkg"
	"io"
	"net/url"
	"os"
)

type Core interface {
	OpKill(
		ctx context.Context,
		opId string,
	)

	NodeCreate()

	NodeKill()

	Run(
		ctx context.Context,
		pkgRef string,
		opts *RunOpts,
	)

	PkgCreate(
		path string,
		description string,
		name string,
	)

	PkgLs(
		path string,
	)

	PkgInstall(
		path,
		pkgRef,
		username,
		password string,
	)

	PkgValidate(
		pkgRef string,
	)

	Events()

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

	opspecNodeURL, err := url.Parse("http://localhost:42224")
	if nil != err {
		panic(err)
	}

	opspecNodeAPIClient := client.New(
		*opspecNodeURL,
		&client.Opts{
			RetryLogHook: func(err error) {
				cliOutput.Attention("request resulted in a recoverable error & will be retried; error was: %v", err.Error())
			},
		},
	)

	return &_core{
		opspecNodeAPIClient:     opspecNodeAPIClient,
		pkg:                     pkg.New(),
		pkgResolver:             newPkgResolver(cliExiter, cliParamSatisfier, *opspecNodeURL),
		cliColorer:              cliColorer,
		cliExiter:               cliExiter,
		cliOutput:               cliOutput,
		cliParamSatisfier:       cliParamSatisfier,
		nodeProvider:            local.New(),
		nodeReachabilityEnsurer: newNodeReachabilityEnsurer(cliExiter),
		updater:                 updater.New(),
		os:                      _os,
		writer:                  os.Stdout,
		ioutil:                  iioutil.New(),
	}

}

type _core struct {
	opspecNodeAPIClient     client.Client
	pkg                     pkg.Pkg
	pkgResolver             pkgResolver
	cliColorer              clicolorer.CliColorer
	cliExiter               cliexiter.CliExiter
	cliOutput               clioutput.CliOutput
	cliParamSatisfier       cliparamsatisfier.CLIParamSatisfier
	nodeProvider            nodeprovider.NodeProvider
	nodeReachabilityEnsurer nodeReachabilityEnsurer
	updater                 updater.Updater
	os                      ios.IOS
	writer                  io.Writer
	ioutil                  iioutil.IIOUtil
}
