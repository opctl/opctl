// Package core defines the core interface for an opspec node
package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/data"
	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/container"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/op"
	"github.com/opctl/opctl/sdks/go/opspec/opfile"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o fakes/core.go . Core
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

	callStore := newCallStore()

	callKiller := newCallKiller(
		callStore,
		containerRuntime,
		pubSub,
	)

	caller := newCaller(
		call.NewInterpreter(
			container.NewInterpreter(dataDirPath),
			dataDirPath,
		),
		newContainerCaller(
			containerRuntime,
			pubSub,
		),
		dataDirPath,
		callStore,
		callKiller,
		pubSub,
	)

	core = _core{
		caller:              caller,
		callKiller:          callKiller,
		containerRuntime:    containerRuntime,
		data:                data.New(),
		dataCachePath:       filepath.Join(dataDirPath, "ops"),
		opFileGetter:        opfile.NewGetter(),
		pubSub:              pubSub,
		uniqueStringFactory: uniqueStringFactory,
	}

	return
}

type _core struct {
	caller              caller
	callKiller          callKiller
	containerRuntime    containerruntime.ContainerRuntime
	data                data.Data
	dataCachePath       string
	opFileGetter        opfile.Getter
	opInterpreter       op.Interpreter
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
}
