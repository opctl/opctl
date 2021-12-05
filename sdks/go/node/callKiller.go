package node

import (
	"context"
	"time"

	"github.com/opctl/opctl/sdks/go/model"

	"github.com/opctl/opctl/sdks/go/node/containerruntime"
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

	for _, childCallGraph := range ckr.stateStore.ListWithParentID(callID) {
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
	}
}
