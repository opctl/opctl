package core

//go:generate counterfeiter -o ./fakeContainerCaller.go --fake-name fakeContainerCaller ./ containerCaller

import (
	"github.com/opspec-io/opctl/pkg/containerengine"
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
	containerEngine containerengine.ContainerEngine,
	pubSub pubsub.PubSub,
	dcgNodeRepo dcgNodeRepo,
) containerCaller {

	return _containerCaller{
		bundle:          bundle,
		containerEngine: containerEngine,
		pubSub:          pubSub,
		dcgNodeRepo:     dcgNodeRepo,
	}

}

type _containerCaller struct {
	bundle          bundle.Bundle
	containerEngine containerengine.ContainerEngine
	pubSub          pubsub.PubSub
	dcgNodeRepo     dcgNodeRepo
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

	containerStartedEvent := model.Event{
		Timestamp: time.Now().UTC(),
		ContainerStarted: &model.ContainerStartedEvent{
			ContainerId: containerId,
			OpRef:       opRef,
			OpGraphId:   opGraphId,
		},
	}
	this.pubSub.Publish(containerStartedEvent)

	err = this.containerEngine.StartContainer(
		newContainerStartReq(inboundScope, scgContainerCall, containerId, op.Inputs, opGraphId),
		this.pubSub,
	)
	if nil != err {
		return
	}

	/* construct outputs */
	outboundScope = map[string]*model.Data{}
	dcgContainer, err := this.containerEngine.InspectContainerIfExists(containerId)
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
	// construct strings
	for scgContainerEnvVarName, scgContainerEnvVar := range scgContainerCall.EnvVars {
		for dcgContainerEnvVarName, dcgContainerEnvVarValue := range dcgContainer.EnvVars {
			if scgContainerEnvVarName == dcgContainerEnvVarName {
				outboundScope[scgContainerEnvVar.Bind] = &model.Data{String: dcgContainerEnvVarValue}
			}
		}
	}
	// construct sockets
	for scgContainerSocketAddress, scgContainerSocket := range scgContainerCall.Sockets {
		for dcgContainerSocketAddress, dcgHostSocketAddress := range dcgContainer.Sockets {
			if scgContainerSocketAddress == dcgContainerSocketAddress {
				outboundScope[scgContainerSocket.Bind] = &model.Data{Socket: dcgHostSocketAddress}
			}
		}
	}

	defer func() {

		this.dcgNodeRepo.DeleteIfExists(containerId)

		this.pubSub.Publish(
			model.Event{
				Timestamp: time.Now().UTC(),
				ContainerExited: &model.ContainerExitedEvent{
					ContainerId: containerId,
					OpRef:       opRef,
					OpGraphId:   opGraphId,
				},
			},
		)

	}()

	return
}
