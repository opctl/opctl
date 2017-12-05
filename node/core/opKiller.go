package core

//go:generate counterfeiter -o ./fakeOpKiller.go --fake-name fakeOpKiller ./ opKiller

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/containerprovider"
	"github.com/opspec-io/sdk-golang/util/pubsub"
)

type opKiller interface {
	Kill(
		rootOpId string,
	)
}

func newOpKiller(
	containerProvider containerprovider.ContainerProvider,
	eventSubscriber pubsub.EventSubscriber,
) opKiller {
	return _opKiller{
		containerProvider: containerProvider,
		eventSubscriber:   eventSubscriber,
	}
}

type _opKiller struct {
	containerProvider containerprovider.ContainerProvider
	eventSubscriber   pubsub.EventSubscriber
}

func (ok _opKiller) Kill(
	rootOpId string,
) {
	for _, containerId := range ok.listContainerIds(rootOpId) {
		err := ok.containerProvider.DeleteContainerIfExists(containerId)
		fmt.Printf(
			"Error encountered killing container w/ id %v, rootOpId %v; error was %v",
			containerId,
			rootOpId,
			err.Error(),
		)
	}
}

func (ok _opKiller) listContainerIds(
	rootOpId string,
) []string {
	eventChannel := make(chan *model.Event, 1)
	ok.eventSubscriber.Subscribe(
		&model.EventFilter{
			Roots: []string{
				rootOpId,
			},
		},
		eventChannel,
	)

	containerIds := []string{}
eventLoop:
	for event := range eventChannel {
		switch {
		case nil != event.OpEnded && event.OpEnded.OpId == rootOpId:
			break eventLoop
		case nil != event.ContainerStarted:
			containerIds = append(containerIds, event.ContainerStarted.ContainerId)
		}
	}
	return containerIds
}
