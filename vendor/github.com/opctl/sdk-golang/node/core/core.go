// Package core defines the core interface for an opspec node
package core

import "github.com/opctl/sdk-golang/model"
import (
	"context"
	"path/filepath"

	"github.com/opctl/sdk-golang/data"
	"github.com/opctl/sdk-golang/node/core/containerruntime"
	"github.com/opctl/sdk-golang/op/dotyml"
	"github.com/opctl/sdk-golang/op/interpreter/containercall"
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

	caller := newCaller(
		newContainerCaller(
			containerRuntime,
			containercall.NewInterpreter(dataDirPath),
			pubSub,
			dcgNodeRepo,
		),
	)

	caller.setParallelCaller(
		newParallelCaller(
			caller,
			opKiller,
			pubSub,
			uniqueStringFactory,
		),
	)

	caller.setSerialCaller(
		newSerialCaller(
			caller,
			pubSub,
			uniqueStringFactory,
		),
	)

	opCaller := newOpCaller(
		pubSub,
		dcgNodeRepo,
		caller,
		dataDirPath,
	)

	caller.setOpCaller(
		opCaller,
	)

	core = _core{
		containerRuntime:    containerRuntime,
		dcgNodeRepo:         dcgNodeRepo,
		opCaller:            opCaller,
		opKiller:            opKiller,
		pubSub:              pubSub,
		dataCachePath:       filepath.Join(dataDirPath, "pkgs"),
		uniqueStringFactory: uniqueStringFactory,
		dotYmlGetter:        dotyml.NewGetter(),
		data:                data.New(),
	}

	return
}

type _core struct {
	containerRuntime    containerruntime.ContainerRuntime
	dcgNodeRepo         dcgNodeRepo
	opCaller            opCaller
	opKiller            opKiller
	pubSub              pubsub.PubSub
	dotYmlGetter        dotyml.Getter
	data                data.Data
	dataCachePath       string
	uniqueStringFactory uniquestring.UniqueStringFactory
}
