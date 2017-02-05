package core

import (
	"github.com/opspec-io/opctl/util/colorer"
	"github.com/opspec-io/opctl/util/updater"
	"github.com/opspec-io/opctl/util/vos"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/engineclient"
	"github.com/opspec-io/sdk-golang/pkg/validate"
	"io"
	"os"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Core

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
	colorer colorer.Colorer,
) Core {

	output := newOutput(colorer, os.Stderr, os.Stdout)
	exiter := newExiter(output, vos.New())

	return &_core{
		bundle:         bundle.New(),
		colorer:        colorer,
		exiter:         exiter,
		engineClient:   engineclient.New(),
		output:         output,
		paramSatisfier: newParamSatisfier(colorer, exiter, output, validate.New(), vos.New()),
		updater:        updater.New(),
		vos:            vos.New(),
		writer:         os.Stdout,
	}

}

type _core struct {
	bundle         bundle.Bundle
	colorer        colorer.Colorer
	exiter         exiter
	engineClient   engineclient.EngineClient
	output         output
	paramSatisfier paramSatisfier
	updater        updater.Updater
	vos            vos.Vos
	writer         io.Writer
}
