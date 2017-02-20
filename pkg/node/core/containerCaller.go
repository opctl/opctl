package core

//go:generate counterfeiter -o ./fakeContainerCaller.go --fake-name fakeContainerCaller ./ containerCaller

import (
	"github.com/opspec-io/opctl/util/containerprovider"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
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
	bundle bundle.Bundle,
	containerProvider containerprovider.ContainerProvider,
	pubSub pubsub.PubSub,
	dcgNodeRepo dcgNodeRepo,
) containerCaller {

	return _containerCaller{
		bundle:            bundle,
		containerProvider: containerProvider,
		pubSub:            pubSub,
		dcgNodeRepo:       dcgNodeRepo,
	}

}

type _containerCaller struct {
	bundle            bundle.Bundle
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

	}()

	op, err := this.bundle.GetOp(
		opRef,
	)
	if nil != err {
		return
	}

	this.dcgNodeRepo.Add(
		&dcgNodeDescriptor{
			Id:        containerId,
			OpRef:     opRef,
			OpGraphId: opGraphId,
			Container: &dcgContainerDescriptor{},
		},
	)

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
		newRunContainerReq(inboundScope, scgContainerCall, containerId, op.Inputs, opGraphId),
		this.pubSub,
	)
	if nil != err {
		return
	}

	/* construct outputs */
	outboundScope = map[string]*model.Data{}
	dcgContainer, err := this.containerProvider.InspectContainerIfExists(containerId)
	// construct files
	for scgContainerFilePath, scgContainerFile := range scgContainerCall.Files {
		for dcgContainerFilePath, dcgHostFilePath := range dcgContainer.Files {
			if scgContainerFilePath == dcgContainerFilePath {
				outboundScope[scgContainerFile.Bind] = &model.Data{File: dcgHostFilePath}
			}
		}
	}
	// construct dirs
	for scgContainerDirPath, scgContainerDir := range scgContainerCall.Dirs {
		for dcgContainerDirPath, dcgHostDirPath := range dcgContainer.Dirs {
			if scgContainerDirPath == dcgContainerDirPath {
				outboundScope[scgContainerDir.Bind] = &model.Data{Dir: dcgHostDirPath}
			}
		}
	}
	// construct sockets
	for scgContainerSocketAddress, scgContainerSocket := range scgContainerCall.Sockets {
		// default; point net sockets @ their containers ip
		outboundScope[scgContainerSocket.Bind] = &model.Data{Socket: dcgContainer.IpAddress}
		for dcgContainerSocketAddress, dcgHostSocketAddress := range dcgContainer.Sockets {
			if scgContainerSocketAddress == dcgContainerSocketAddress {
				// override default; point unix sockets @ their location on the host
				outboundScope[scgContainerSocket.Bind] = &model.Data{Socket: dcgHostSocketAddress}
			}
		}
	}

	return
}
