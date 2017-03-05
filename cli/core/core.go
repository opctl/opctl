package core

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Core

import (
	"github.com/opspec-io/opctl/pkg/nodeprovider"
	"github.com/opspec-io/opctl/pkg/nodeprovider/local"
	"github.com/opspec-io/opctl/util/clicolorer"
	"github.com/opspec-io/opctl/util/cliexiter"
	"github.com/opspec-io/opctl/util/clioutput"
	"github.com/opspec-io/opctl/util/cliparamsatisfier"
	"github.com/opspec-io/opctl/util/updater"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/consumenodeapi"
	"github.com/opspec-io/sdk-golang/pkg/pkg"
	"github.com/opspec-io/sdk-golang/pkg/validate"
	"io"
	"os"
)

type Core interface {
	CreateCollection(
		description string,
		name string,
	)

	CreateOp(
		collection string,
		description string,
		name string,
	)

	KillOp(
		opId string,
	)

	ListOpsInCollection(
		collection string,
	)

	NodeCreate()

	NodeKill()

	RunOp(
		args []string,
		collection string,
		name string,
	)

	SetCollectionDescription(
		description string,
	)

	SetOpDescription(
		collection string,
		description string,
		name string,
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
