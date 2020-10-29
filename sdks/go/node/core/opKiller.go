package core

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/opctl/opctl/sdks/go/model"

	"github.com/opctl/opctl/sdks/go/node/core/containerruntime"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

//counterfeiter:generate -o internal/fakes/opKiller.go . opKiller
type opKiller interface {
	Kill(
		callID string,
		rootCallID string,
	)
}

func newOpKiller(
	stateStore stateStore,
	containerRuntime containerruntime.ContainerRuntime,
	eventPublisher pubsub.EventPublisher,
) opKiller {
	return _opKiller{
		stateStore:       stateStore,
		containerRuntime: containerRuntime,
		eventPublisher:   eventPublisher,
	}
}

type _opKiller struct {
	stateStore       stateStore
	containerRuntime containerruntime.ContainerRuntime
	eventPublisher   pubsub.EventPublisher
}

func (ckr _opKiller) Kill(
	callID string,
	rootCallID string,
) {
	ckr.containerRuntime.DeleteContainerIfExists(callID)

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
					OpKillRequested: &model.OpKillRequested{
						Request: model.KillOpReq{
							OpID:       childCallGraph.Id,
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
