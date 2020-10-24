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
	callStore callStore,
	containerRuntime containerruntime.ContainerRuntime,
	eventPublisher pubsub.EventPublisher,
) opKiller {
	return _opKiller{
		callStore:        callStore,
		containerRuntime: containerRuntime,
		eventPublisher:   eventPublisher,
	}
}

type _opKiller struct {
	callStore        callStore
	containerRuntime containerruntime.ContainerRuntime
	eventPublisher   pubsub.EventPublisher
}

func (ckr _opKiller) Kill(
	callID string,
	rootCallID string,
) {
	ckr.eventPublisher.Publish(
		model.Event{
			OpKillRequested: &model.OpKillRequested{
				Request: model.KillOpReq{
					OpID:     callID,
					RootOpID: rootCallID,
				},
			},
			Timestamp: time.Now().UTC(),
		},
	)

	ckr.callStore.SetIsKilled(callID)
	ckr.containerRuntime.DeleteContainerIfExists(callID)

	var waitGroup sync.WaitGroup

	for _, childCallGraph := range ckr.callStore.ListWithParentID(callID) {
		// recursively kill all child calls
		waitGroup.Add(1)
		go func(childCallGraph *model.DCG) {
			defer func() {
				if panicArg := recover(); panicArg != nil {
					// recover from panics; treat as errors
					fmt.Println(panicArg, debug.Stack())
				}
			}()
			defer waitGroup.Done()

			ckr.Kill(
				childCallGraph.Id,
				rootCallID,
			)

		}(childCallGraph)
	}

	waitGroup.Wait()

}
