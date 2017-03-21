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
	"github.com/opctl/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/consumenodeapi"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/validate"
	"io"
	"os"
)

type Core interface {
	Create(
		path string,
		description string,
		name string,
	)

	OpKill(
		opId string,
	)

	ListPackages(
		path string,
	)

	NodeCreate()

	NodeKill()

	RunOp(
		args []string,
		pkgRef string,
	)

	PkgSetDescription(
		description string,
		pkgRef string,
	)

	PkgValidate(
		pkgRef string,
	)

	StreamEvents()

	SelfUpdate(
		channel string,
	)
}

func New(
	cliColorer clicolorer.CliColorer,
) Core {

	cliOutput := clioutput.New(cliColorer, os.Stderr, os.Stdout)
	cliExiter := cliexiter.New(cliOutput, vos.New())

	return &_core{
		consumeNodeApi:    consumenodeapi.New(),
		pkg:               pkg.New(),
		cliColorer:        cliColorer,
		cliExiter:         cliExiter,
		cliOutput:         cliOutput,
		cliParamSatisfier: cliparamsatisfier.New(cliColorer, cliExiter, cliOutput, validate.New(), vos.New()),
		nodeProvider:      local.New(),
		updater:           updater.New(),
		vos:               vos.New(),
		writer:            os.Stdout,
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
	vos               vos.Vos
	writer            io.Writer
}
