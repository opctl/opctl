package core

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Core

import (
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
	"github.com/opspec-io/sdk-golang/validate"
	"io"
	"net/url"
	"os"
)

type Core interface {
	OpKill(
		opId string,
	)

	NodeCreate()

	NodeKill()

	Run(
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

	PkgPull(
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

	opspecNodeURL, err := url.Parse("http://localhost:42224")
	if nil != err {
		panic(err)
	}

	return &_core{
		opspecNodeAPIClient: client.New(*opspecNodeURL),
		pkg:                 pkg.New(),
		cliColorer:          cliColorer,
		cliExiter:           cliExiter,
		cliOutput:           cliOutput,
		cliParamSatisfier:   cliparamsatisfier.New(cliExiter, cliOutput, validate.New()),
		nodeProvider:        local.New(),
		updater:             updater.New(),
		os:                  _os,
		writer:              os.Stdout,
		ioutil:              iioutil.New(),
	}

}

type _core struct {
	opspecNodeAPIClient client.Client
	pkg                 pkg.Pkg
	cliColorer          clicolorer.CliColorer
	cliExiter           cliexiter.CliExiter
	cliOutput           clioutput.CliOutput
	cliParamSatisfier   cliparamsatisfier.CliParamSatisfier
	nodeProvider        nodeprovider.NodeProvider
	updater             updater.Updater
	os                  ios.IOS
	writer              io.Writer
	ioutil              iioutil.Iioutil
}
