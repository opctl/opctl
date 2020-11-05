package core

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/opctl/opctl/sdks/go/model"

	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o internal/fakes/callKiller.go . callKiller
type callKiller interface {
	Kill(
		ctx context.Context,
		callID string,
		rootCallID string,
	)
}

func newCallKiller(
	stateStore stateStore,
	containerRuntime containerruntime.ContainerRuntime,
	eventPublisher pubsub.EventPublisher,
) callKiller {
	return _callKiller{
		stateStore:       stateStore,
		containerRuntime: containerRuntime,
		eventPublisher:   eventPublisher,
	}
}

type _callKiller struct {
	stateStore       stateStore
	containerRuntime containerruntime.ContainerRuntime
	eventPublisher   pubsub.EventPublisher
}

func (ckr _callKiller) Kill(
	ctx context.Context,
	callID string,
	rootCallID string,
) {
	ckr.containerRuntime.DeleteContainerIfExists(
		ctx,
		callID,
	)

	var waitGroup sync.WaitGroup

	for _, childCallGraph := range ckr.stateStore.ListWithParentID(callID) {
		// recursively kill all child calls
		waitGroup.Add(1)
		go func(childCallGraph *model.Call) {
			defer func() {
				if panicArg := recover(); panicArg != nil {
					// recover from panics; treat as errors
					fmt.Println(panicArg, debug.Stack())
				}
			}()
			defer waitGroup.Done()

			ckr.eventPublisher.Publish(
				model.Event{
					CallKillRequested: &model.CallKillRequested{
						Request: model.KillOpReq{
							OpID:       childCallGraph.ID,
							RootCallID: rootCallID,
						},
					},
					Timestamp: time.Now().UTC(),
				},
			)

		}(childCallGraph)
	}

	waitGroup.Wait()

}
