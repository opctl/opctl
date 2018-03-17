// Package core defines the core interface for an opspec node
package core

import "github.com/opspec-io/sdk-golang/model"
import (
	"context"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/node/core/containerruntime"
	"github.com/opspec-io/sdk-golang/op/interpreter/containercall"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
	"path/filepath"
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

	// Resolve attempts to resolve a pkg via local filesystem or git
	// nil pullCreds will be ignored
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	ResolvePkg(
		ctx context.Context,
		pkgRef string,
		pullCreds *model.PullCreds,
	) (
		model.DataHandle,
		error,
	)
}

func New(
	pubSub pubsub.PubSub,
	containerRuntime containerruntime.ContainerRuntime,
	rootFSPath string,
) (core Core) {
	uniqueStringFactory := uniquestring.NewUniqueStringFactory()

	dcgNodeRepo := newDCGNodeRepo()

	opKiller := newOpKiller(dcgNodeRepo, containerRuntime)

	caller := newCaller(
		newContainerCaller(
			containerRuntime,
			containercall.NewInterpreter(rootFSPath),
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
		rootFSPath,
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
		pkgCachePath:        filepath.Join(rootFSPath, "pkgs"),
		uniqueStringFactory: uniqueStringFactory,
		pkg:                 pkg.New(),
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
	pkg                 pkg.Pkg
	data                data.Data
	pkgCachePath        string
	uniqueStringFactory uniquestring.UniqueStringFactory
}
