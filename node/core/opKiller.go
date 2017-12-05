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
	containerIdChan := make(chan string, 1)
	ok.listContainerIds(rootOpId, containerIdChan)

	for containerId := range containerIdChan {
		go func(containerId string) {
			if err := ok.containerProvider.DeleteContainerIfExists(containerId); nil != err {
				fmt.Printf(
					"Error encountered killing container w/ id %v, rootOpId %v; error was %v",
					containerId,
					rootOpId,
					err.Error(),
				)
			}
		}(containerId)
	}
}

func (ok _opKiller) listContainerIds(
	rootOpId string,
	containerIdChannel chan string,
) {
	eventChannel := make(chan *model.Event, 1)
	ok.eventSubscriber.Subscribe(
		&model.EventFilter{
			Roots: []string{
				rootOpId,
			},
		},
		eventChannel,
	)

	for event := range eventChannel {
		switch {
		case nil != event.ContainerStarted:
			containerIdChannel <- event.ContainerStarted.ContainerId
		case nil != event.OpKilled && rootOpId == event.OpKilled.RootOpId:
			close(containerIdChannel)
			return
		}
	}
}
