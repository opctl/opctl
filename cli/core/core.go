package core

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Core

import (
	"github.com/opctl/opctl/nodeprovider"
	"github.com/opctl/opctl/nodeprovider/local"
	"github.com/opctl/opctl/util/clicolorer"
	"github.com/opctl/opctl/util/cliexiter"
	"github.com/opctl/opctl/util/clioutput"
	"github.com/opctl/opctl/util/cliparamsatisfier"
	"github.com/opctl/opctl/util/updater"
	"github.com/opspec-io/sdk-golang/consumenodeapi"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/validate"
	"github.com/virtual-go/fs/osfs"
	"github.com/virtual-go/vioutil"
	"github.com/virtual-go/vos"
	"io"
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

	_fs := osfs.New()
	_os := vos.New(_fs)

	cliOutput := clioutput.New(cliColorer, os.Stderr, os.Stdout)
	cliExiter := cliexiter.New(cliOutput, _os)

	return &_core{
		consumeNodeApi:    consumenodeapi.New(),
		pkg:               pkg.New(),
		cliColorer:        cliColorer,
		cliExiter:         cliExiter,
		cliOutput:         cliOutput,
		cliParamSatisfier: cliparamsatisfier.New(cliExiter, cliOutput, validate.New()),
		nodeProvider:      local.New(),
		updater:           updater.New(),
		os:                _os,
		writer:            os.Stdout,
		ioutil:            vioutil.New(_fs),
	}

}

type _core struct {
	consumeNodeApi    consumenodeapi.ConsumeNodeApi
	pkg               pkg.Pkg
	cliColorer        clicolorer.CliColorer
	cliExiter         cliexiter.CliExiter
	cliOutput         clioutput.CliOutput
	cliParamSatisfier cliparamsatisfier.CliParamSatisfier
	nodeProvider      nodeprovider.NodeProvider
	updater           updater.Updater
	os                vos.VOS
	writer            io.Writer
	ioutil            vioutil.VIOUtil
}
