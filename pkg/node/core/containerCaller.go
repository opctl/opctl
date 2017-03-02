package core

//go:generate counterfeiter -o ./fakeContainerCaller.go --fake-name fakeContainerCaller ./ containerCaller

import (
	"github.com/opspec-io/opctl/util/containerprovider"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"strings"
	"time"
)

type containerCaller interface {
	// Executes a container call
	Call(
		inboundScope map[string]*model.Data,
		containerId string,
		scgContainerCall *model.ScgContainerCall,
		opRef string,
		opGraphId string,
	) (
		outboundScope map[string]*model.Data,
		err error,
	)
}

func newContainerCaller(
	containerProvider containerprovider.ContainerProvider,
	pubSub pubsub.PubSub,
	dcgNodeRepo dcgNodeRepo,
) containerCaller {

	return _containerCaller{
		containerProvider: containerProvider,
		pubSub:            pubSub,
		dcgNodeRepo:       dcgNodeRepo,
	}

}

type _containerCaller struct {
	containerProvider containerprovider.ContainerProvider
	pubSub            pubsub.PubSub
	dcgNodeRepo       dcgNodeRepo
}

func (this _containerCaller) Call(
	inboundScope map[string]*model.Data,
	containerId string,
	scgContainerCall *model.ScgContainerCall,
	opRef string,
	opGraphId string,
) (
	outboundScope map[string]*model.Data,
	err error,
) {
	defer func() {
		// defer must be defined before conditional return statements so it always runs

		this.dcgNodeRepo.DeleteIfExists(containerId)

	}()

	this.dcgNodeRepo.Add(
		&dcgNodeDescriptor{
			Id:        containerId,
			OpRef:     opRef,
			OpGraphId: opGraphId,
			Container: &dcgContainerDescriptor{},
		},
	)

	runContainerReq, err := newRunContainerReq(inboundScope, scgContainerCall, containerId, opGraphId, opRef)
	if nil != err {
		return
	}

	this.pubSub.Publish(
		&model.Event{
			Timestamp: time.Now().UTC(),
			ContainerStarted: &model.ContainerStartedEvent{
				ContainerId: containerId,
				OpRef:       opRef,
				OpGraphId:   opGraphId,
			},
		},
	)
	err = this.containerProvider.RunContainer(
		runContainerReq,
		this.pubSub,
	)
	this.pubSub.Publish(
		&model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExitedEvent{
				ContainerId: containerId,
				OpRef:       opRef,
				OpGraphId:   opGraphId,
			},
		},
	)
	if nil != err {
		return
	}

	/* construct outputs */
	outboundScope, err = this.constructOutboundScope(containerId, opGraphId, scgContainerCall)

	return
}

// O(n) complexity (n being events in op graph)
func (this _containerCaller) constructOutboundScope(
	containerId string,
	opGraphId string,
	scgContainerCall *model.ScgContainerCall,
) (outboundScope map[string]*model.Data, err error) {
	outboundScope = map[string]*model.Data{}

	// subscribe to op graph events
	eventChannel := make(chan *model.Event)
	this.pubSub.Subscribe(
		&model.EventFilter{OpGraphIds: []string{opGraphId}},
		eventChannel,
	)

eventLoop:
	for event := range eventChannel {
		switch {
		case nil != event.ContainerExited && event.ContainerExited.ContainerId == containerId:
			break eventLoop
		case nil != event.ContainerStdErrWrittenTo && event.ContainerStdErrWrittenTo.ContainerId == containerId:
			for boundPrefix, binding := range scgContainerCall.StdErr {
				rawOutput := string(event.ContainerStdErrWrittenTo.Data)
				trimmedOutput := strings.TrimPrefix(rawOutput, boundPrefix)
				if trimmedOutput != rawOutput {
					// if output trimming had effect we've got a match
					outboundScope[binding] = &model.Data{String: trimmedOutput}
				}
			}
		case nil != event.ContainerStdOutWrittenTo && event.ContainerStdOutWrittenTo.ContainerId == containerId:
			for boundPrefix, binding := range scgContainerCall.StdOut {
				rawOutput := string(event.ContainerStdOutWrittenTo.Data)
				trimmedOutput := strings.TrimPrefix(rawOutput, boundPrefix)
				if trimmedOutput != rawOutput {
					// if output trimming had effect we've got a match
					outboundScope[binding] = &model.Data{String: trimmedOutput}
				}
			}
		}
	}

	dcgContainer, err := this.containerProvider.InspectContainerIfExists(containerId)
	// construct files
	for scgContainerFilePath, scgContainerFile := range scgContainerCall.Files {
		for dcgContainerFilePath, dcgHostFilePath := range dcgContainer.Files {
			if scgContainerFilePath == dcgContainerFilePath {
				outboundScope[scgContainerFile] = &model.Data{File: dcgHostFilePath}
			}
		}
	}
	// construct dirs
	for scgContainerDirPath, scgContainerDir := range scgContainerCall.Dirs {
		for dcgContainerDirPath, dcgHostDirPath := range dcgContainer.Dirs {
			if scgContainerDirPath == dcgContainerDirPath {
				outboundScope[scgContainerDir] = &model.Data{Dir: dcgHostDirPath}
			}
		}
	}
	// construct sockets
	for scgContainerSocketAddress, scgContainerSocket := range scgContainerCall.Sockets {
		// default; point net sockets @ their containers ip
		outboundScope[scgContainerSocket] = &model.Data{Socket: dcgContainer.IpAddress}
		for dcgContainerSocketAddress, dcgHostSocketAddress := range dcgContainer.Sockets {
			if scgContainerSocketAddress == dcgContainerSocketAddress {
				// override default; point unix sockets @ their location on the host
				outboundScope[scgContainerSocket] = &model.Data{Socket: dcgHostSocketAddress}
			}
		}
	}

	return
}
