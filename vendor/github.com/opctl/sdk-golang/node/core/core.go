// Package core defines the core interface for an opspec node
package core

import (
	"context"
	"path/filepath"

	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/core/containerruntime"
	"github.com/opctl/sdk-golang/opspec/interpreter/call"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/container"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/op"
	"github.com/opctl/sdk-golang/opspec/opfile"
	"github.com/opctl/sdk-golang/util/pubsub"
	"github.com/opctl/sdk-golang/util/uniquestring"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Core

type Core interface {
	GetEventStream(
		ctx context.Context,
		req *model.GetEventStreamReq,
	) (
		<-chan model.Event,
		<-chan error,
	)

	KillOp(
		req model.KillOpReq,
	)

	StartOp(
		ctx context.Context,
		req model.StartOpReq,
	) (
		callId string,
		err error,
	)

	// Resolve attempts to resolve an op via local filesystem or git
	// nil pullCreds will be ignored
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	ResolveData(
		ctx context.Context,
		dataRef string,
		pullCreds *model.PullCreds,
	) (
		model.DataHandle,
		error,
	)
}

func New(
	pubSub pubsub.PubSub,
	containerRuntime containerruntime.ContainerRuntime,
	dataDirPath string,
) (core Core) {
	uniqueStringFactory := uniquestring.NewUniqueStringFactory()

	dcgNodeRepo := newDCGNodeRepo()

	opKiller := newOpKiller(dcgNodeRepo, containerRuntime)

	opInterpreter := op.NewInterpreter(dataDirPath)

	caller := newCaller(
		call.NewInterpreter(
			container.NewInterpreter(dataDirPath),
			dataDirPath,
		),
		newContainerCaller(
			containerRuntime,
			pubSub,
			dcgNodeRepo,
		),
		dataDirPath,
		dcgNodeRepo,
		opKiller,
		pubSub,
	)

	core = _core{
		caller:              caller,
		containerRuntime:    containerRuntime,
		data:                data.New(),
		dataCachePath:       filepath.Join(dataDirPath, "pkgs"),
		dcgNodeRepo:         dcgNodeRepo,
		dotYmlGetter:        dotyml.NewGetter(),
		opInterpreter:       opInterpreter,
		opKiller:            opKiller,
		pubSub:              pubSub,
		uniqueStringFactory: uniqueStringFactory,
	}

	return
}

type _core struct {
	caller              caller
	containerRuntime    containerruntime.ContainerRuntime
	data                data.Data
	dataCachePath       string
	dcgNodeRepo         dcgNodeRepo
	dotYmlGetter        dotyml.Getter
	opInterpreter       op.Interpreter
	opKiller            opKiller
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
}
