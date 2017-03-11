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
		scope map[string]*model.Data,
		outputs chan *variable,
		containerId string,
		scgContainerCall *model.ScgContainerCall,
		pkgRef string,
		rootOpId string,
	) (
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
	scope map[string]*model.Data,
	outputs chan *variable,
	containerId string,
	scgContainerCall *model.ScgContainerCall,
	pkgRef string,
	rootOpId string,
) (
	err error,
) {
	defer func() {
		// defer must be defined before conditional return statements so it always runs

		this.dcgNodeRepo.DeleteIfExists(containerId)

		this.containerProvider.DeleteContainerIfExists(containerId)

	}()

	this.dcgNodeRepo.Add(
		&dcgNodeDescriptor{
			Id:        containerId,
			PkgRef:    pkgRef,
			RootOpId:  rootOpId,
			Container: &dcgContainerDescriptor{},
		},
	)

	dcgContainerCall, err := constructDcgContainerCall(scope, scgContainerCall, containerId, rootOpId, pkgRef)
	if nil != err {
		return
	}

	go this.sendOutputs(dcgContainerCall, outputs, scgContainerCall)

	this.pubSub.Publish(
		&model.Event{
			Timestamp: time.Now().UTC(),
			ContainerStarted: &model.ContainerStartedEvent{
				ContainerId: containerId,
				PkgRef:      pkgRef,
				RootOpId:    rootOpId,
			},
		},
	)
	err = this.containerProvider.RunContainer(
		dcgContainerCall,
		this.pubSub,
	)
	this.pubSub.Publish(
		&model.Event{
			Timestamp: time.Now().UTC(),
			ContainerExited: &model.ContainerExitedEvent{
				ContainerId: containerId,
				PkgRef:      pkgRef,
				RootOpId:    rootOpId,
			},
		},
	)
	close(outputs)

	return
}

func (this _containerCaller) sendOutputs(
	dcgContainerCall *model.DcgContainerCall,
	outputs chan *variable,
	scgContainerCall *model.ScgContainerCall,
) {

	// send socket outputs
	for socketAddr, varName := range scgContainerCall.Sockets {
		if "0.0.0.0" == socketAddr {
			outputs <- &variable{
				Name:  varName,
				Value: &model.Data{Socket: dcgContainerCall.ContainerId},
			}
		}
	}

	// send file outputs
	for scgContainerFilePath, varName := range scgContainerCall.Files {
		for dcgContainerFilePath, dcgHostFilePath := range dcgContainerCall.Files {
			if scgContainerFilePath == dcgContainerFilePath {
				outputs <- &variable{
					Name:  varName,
					Value: &model.Data{File: dcgHostFilePath},
				}
			}
		}
	}

	// send dir outputs
	for scgContainerDirPath, varName := range scgContainerCall.Dirs {
		for dcgContainerDirPath, dcgHostDirPath := range dcgContainerCall.Dirs {
			if scgContainerDirPath == dcgContainerDirPath {
				outputs <- &variable{
					Name:  varName,
					Value: &model.Data{Dir: dcgHostDirPath},
				}
			}
		}
	}

	// subscribe to op graph events
	eventChannel := make(chan *model.Event)
	this.pubSub.Subscribe(
		&model.EventFilter{RootOpIds: []string{dcgContainerCall.RootOpId}},
		eventChannel,
	)

	// send string outputs
eventLoop:
	for event := range eventChannel {
		switch {
		case nil != event.ContainerExited && event.ContainerExited.ContainerId == dcgContainerCall.ContainerId:
			break eventLoop
		case nil != event.ContainerStdErrWrittenTo &&
			event.ContainerStdErrWrittenTo.ContainerId == dcgContainerCall.ContainerId:
			for boundPrefix, varName := range scgContainerCall.StdErr {
				rawOutput := string(event.ContainerStdErrWrittenTo.Data)
				trimmedOutput := strings.TrimPrefix(rawOutput, boundPrefix)
				if trimmedOutput != rawOutput {
					// if output trimming had effect we've got a match
					outputs <- &variable{
						Name:  varName,
						Value: &model.Data{String: trimmedOutput},
					}
				}
			}
		case nil != event.ContainerStdOutWrittenTo &&
			event.ContainerStdOutWrittenTo.ContainerId == dcgContainerCall.ContainerId:
			for boundPrefix, varName := range scgContainerCall.StdOut {
				rawOutput := string(event.ContainerStdOutWrittenTo.Data)
				trimmedOutput := strings.TrimPrefix(rawOutput, boundPrefix)
				if trimmedOutput != rawOutput {
					// if output trimming had effect we've got a match
					outputs <- &variable{
						Name:  varName,
						Value: &model.Data{String: trimmedOutput},
					}
				}
			}
		}
	}

	return
}
