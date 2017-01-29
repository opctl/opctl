package core

//go:generate counterfeiter -o ./fakeContainerCaller.go --fake-name fakeContainerCaller ./ containerCaller

import (
	"fmt"
	"github.com/opspec-io/opctl/pkg/containerengine"
	"github.com/opspec-io/opctl/util/eventbus"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"time"
)

type containerCaller interface {
	// Executes a container call
	Call(
		args map[string]*model.Data,
		containerId string,
		containerCall *model.ScgContainerCall,
		opRef string,
		opGraphId string,
	) (
		outputs map[string]*model.Data,
		err error,
	)
}

func newContainerCaller(
	bundle bundle.Bundle,
	containerEngine containerengine.ContainerEngine,
	eventBus eventbus.EventBus,
	dcgNodeRepo dcgNodeRepo,
) containerCaller {

	return &_containerCaller{
		bundle:          bundle,
		containerEngine: containerEngine,
		eventBus:        eventBus,
		dcgNodeRepo:     dcgNodeRepo,
	}

}

type _containerCaller struct {
	bundle          bundle.Bundle
	containerEngine containerengine.ContainerEngine
	eventBus        eventbus.EventBus
	dcgNodeRepo     dcgNodeRepo
}

func (this _containerCaller) Call(
	args map[string]*model.Data,
	containerId string,
	containerCall *model.ScgContainerCall,
	opRef string,
	opGraphId string,
) (
	outputs map[string]*model.Data,
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
	this.eventBus.Publish(containerStartedEvent)

	err = this.containerEngine.StartContainer(
		newContainerStartReq(args, containerCall, containerId, op.Inputs, opGraphId),
		this.eventBus,
	)
	if nil != err {
		return
	}

	// construct outputs
	outputs = map[string]*model.Data{}
	container, err := this.containerEngine.InspectContainerIfExists(containerId)
	fmt.Println(opRef)
	fmt.Printf("containerCaller.container:\n %#v\n", container)
	// construct files
	for scgContainerFilePath, scgContainerFile := range containerCall.Files {
		for containerFilePath, hostFilePath := range container.Files {
			if scgContainerFilePath == containerFilePath {
				outputs[scgContainerFile.Bind] = &model.Data{File: hostFilePath}
			}
		}
	}
	// construct dirs
	for scgContainerDirPath, scgContainerDir := range containerCall.Dirs {
		for containerDirPath, hostDirPath := range container.Dirs {
			if scgContainerDirPath == containerDirPath {
				outputs[scgContainerDir.Bind] = &model.Data{Dir: hostDirPath}
			}
		}
	}
	// construct strings
	for scgContainerEnvVarName, scgContainerEnvVar := range containerCall.EnvVars {
		for containerEnvVarName, containerEnvVarValue := range container.EnvVars {
			if scgContainerEnvVarName == containerEnvVarName {
				outputs[scgContainerEnvVar.Bind] = &model.Data{String: containerEnvVarValue}
			}
		}
	}
	// construct sockets
	for scgContainerSocketAddress, scgContainerSocket := range containerCall.Sockets {
		for containerSocketAddress, hostSocketAddress := range container.Sockets {
			if scgContainerSocketAddress == containerSocketAddress {
				outputs[scgContainerSocket.Bind] = &model.Data{Socket: hostSocketAddress}
			}
		}
	}
	fmt.Printf("containerCaller.outputs:\n %#v\n", outputs)

	defer func() {

		this.dcgNodeRepo.DeleteIfExists(containerId)

		this.eventBus.Publish(
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
