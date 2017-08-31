// Package core defines the core interface for an opspec node
package core

import "github.com/opspec-io/sdk-golang/model"
import (
	"github.com/opspec-io/sdk-golang/containercall"
	"github.com/opspec-io/sdk-golang/pkg"
	"github.com/opspec-io/sdk-golang/util/containerprovider"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"github.com/opspec-io/sdk-golang/util/uniquestring"
	"path/filepath"
)

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Core

type Core interface {
	GetEventStream(
		req *model.GetEventStreamReq,
		eventChannel chan *model.Event,
	) error

	KillOp(
		req model.KillOpReq,
	)

	StartOp(
		req model.StartOpReq,
	) (
		callId string,
		err error,
	)

	// Resolve attempts to resolve a pkg via local filesystem or git
	// nil pullCreds will be ignored
  //
  // expected errs:
  //  - ErrPkgPullAuthentication on authentication failure
  //  - ErrPkgPullAuthorization on authorization failure
  //  - ErrPkgNotFound on resolution failure
	ResolvePkg(
		pkgRef string,
		pullCreds *model.PullCreds,
	) (
		model.PkgHandle,
		error,
	)
}

func New(
	pubSub pubsub.PubSub,
	containerProvider containerprovider.ContainerProvider,
	rootFSPath string,
) (core Core) {
	uniqueStringFactory := uniquestring.NewUniqueStringFactory()

	dcgNodeRepo := newDCGNodeRepo()

	opKiller := newOpKiller(dcgNodeRepo, containerProvider)

	caller := newCaller(
		newContainerCaller(
			containerProvider,
			containercall.New(rootFSPath),
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
		containerProvider:   containerProvider,
		dcgNodeRepo:         dcgNodeRepo,
		opCaller:            opCaller,
		opKiller:            opKiller,
		pubSub:              pubSub,
		pkgCachePath:        filepath.Join(rootFSPath, "pkgs"),
		uniqueStringFactory: uniqueStringFactory,
		pkg:                 pkg.New(),
	}

	return
}

type _core struct {
	containerProvider   containerprovider.ContainerProvider
	dcgNodeRepo         dcgNodeRepo
	opCaller            opCaller
	opKiller            opKiller
	pubSub              pubsub.PubSub
	pkg                 pkg.Pkg
	pkgCachePath        string
	uniqueStringFactory uniquestring.UniqueStringFactory
}
