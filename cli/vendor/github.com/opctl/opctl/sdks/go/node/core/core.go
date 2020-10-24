// Package core defines the core interface for an opspec node
package core

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"path/filepath"
	"time"

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
) Core {
	uniqueStringFactory := uniquestring.NewUniqueStringFactory()

	callStore := newCallStore()

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
		pubSub,
	)

	go func() {
		// process events in background
		opKiller := newOpKiller(
			callStore,
			containerRuntime,
			pubSub,
		)

		since := time.Now().UTC()
		eventChannel, _ := pubSub.Subscribe(
			context.Background(),
			model.EventFilter{
				Since: &since,
			},
		)

		for event := range eventChannel {
			switch {
			case nil != event.OpKillRequested:
				req := event.OpKillRequested.Request
				opKiller.Kill(req.OpID, req.OpID)
			}
		}
	}()

	return _core{
		caller:           caller,
		containerRuntime: containerRuntime,
		data:             data.New(),
		dataCachePath:    filepath.Join(dataDirPath, "ops"),
		opCaller: newOpCaller(
			callStore,
			pubSub,
			caller,
			dataDirPath,
		),
		opFileGetter:        opfile.NewGetter(),
		opInterpreter:       op.NewInterpreter(dataDirPath),
		pubSub:              pubSub,
		uniqueStringFactory: uniqueStringFactory,
	}
}

type _core struct {
	caller              caller
	containerRuntime    containerruntime.ContainerRuntime
	data                data.Data
	dataCachePath       string
	opCaller            opCaller
	opFileGetter        opfile.Getter
	opInterpreter       op.Interpreter
	pubSub              pubsub.PubSub
	uniqueStringFactory uniquestring.UniqueStringFactory
}
