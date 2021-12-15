// Package node defines the core interface for an opctl node
package node

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

// New returns a new LocalCore initialized with the given options
func New(
	ctx context.Context,
	containerRuntime containerruntime.ContainerRuntime,
	dataDirPath string,
) Core {
	eventDbPath := path.Join(dataDirPath, "dcg", "events")
	err := os.MkdirAll(eventDbPath, 0700)
	if err != nil {
		panic(err)
	}

	// per badger README.MD#FAQ "maximizes throughput"
	runtime.GOMAXPROCS(128)

	db, err := badger.Open(
		badger.DefaultOptions(
			eventDbPath,
		).WithLogger(nil),
	)
	if err != nil {
		panic(err)
	}

	pubSub := pubsub.New(db)

	stateStore := newStateStore(
		ctx,
		db,
		pubSub,
	)

	caller := newCaller(
		newContainerCaller(
			containerRuntime,
			pubSub,
			stateStore,
		),
		dataDirPath,
		pubSub,
	)

	go func() {
		// process events in background
		callKiller := newCallKiller(
			stateStore,
			containerRuntime,
			pubSub,
		)

		since := time.Now().UTC()
		eventChannel, _ := pubSub.Subscribe(
			ctx,
			model.EventFilter{
				Since: &since,
			},
		)

		for event := range eventChannel {
			switch {
			case event.CallKillRequested != nil:
				req := event.CallKillRequested.Request
				callKiller.Kill(
					ctx,
					req.OpID,
					req.RootCallID,
				)
			}
		}
	}()

	return core{
		caller:           caller,
		containerRuntime: containerRuntime,
		dataCachePath:    filepath.Join(dataDirPath, "ops"),
		opCaller: newOpCaller(
			caller,
			dataDirPath,
		),
		pubSub:     pubSub,
		stateStore: stateStore,
	}
}

// core is an Node that supports running ops directly on the host
type core struct {
	caller           caller
	containerRuntime containerruntime.ContainerRuntime
	dataCachePath    string
	opCaller         opCaller
	pubSub           pubsub.PubSub
	stateStore       stateStore
}

func (c core) Liveness(
	ctx context.Context,
) error {
	return nil
}

//counterfeiter:generate -o fakes/core.go . Core

// Core is an Node that supports running ops directly on the current machine
type Core interface {
	Node

	// Resolve attempts to resolve data via local filesystem or git
	// nil pullCreds will be ignored
	//
	// expected errs:
	//  - ErrDataProviderAuthentication on authentication failure
	//  - ErrDataProviderAuthorization on authorization failure
	//  - ErrDataRefResolution on resolution failure
	ResolveData(
		ctx context.Context,
		dataRef string,
		pullCreds *model.Creds,
	) (
		model.DataHandle,
		error,
	)
}
