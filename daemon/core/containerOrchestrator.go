package core

import (
	"fmt"
	"github.com/opspec-io/opctl/pkg/containerengine"
	"github.com/opspec-io/opctl/util/eventbus"
	"github.com/opspec-io/opctl/util/interpolater"
	"github.com/opspec-io/sdk-golang/pkg/bundle"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"time"
)

type containerOrchestrator interface {
	Execute(
		args map[string]*model.Arg,
		containerId string,
		containerCall *model.ContainerCall,
		opRef string,
		opGraphId string,
	) (err error)
}

func newContainerOrchestrator(
	bundle bundle.Bundle,
	containerEngine containerengine.ContainerEngine,
	eventBus eventbus.EventBus,
	nodeRepo nodeRepo,
) containerOrchestrator {

	return &_containerOrchestrator{
		bundle:          bundle,
		containerEngine: containerEngine,
		eventBus:        eventBus,
		nodeRepo:        nodeRepo,
	}

}

type _containerOrchestrator struct {
	bundle          bundle.Bundle
	containerEngine containerengine.ContainerEngine
	eventBus        eventbus.EventBus
	nodeRepo        nodeRepo
}

func (this _containerOrchestrator) Execute(
	args map[string]*model.Arg,
	containerId string,
	containerCall *model.ContainerCall,
	opRef string,
	opGraphId string,
) (err error) {

	op, err := this.bundle.GetOp(
		opRef,
	)
	if nil != err {
		return
	}

	this.nodeRepo.add(
		&nodeDescriptor{
			Id:        containerId,
			OpRef:     opRef,
			OpGraphId: opGraphId,
			Container: &containerDescriptor{},
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

	// build cmd, env, fs, & net
	cmd := []string{}
	env := []*model.ContainerInstanceEnvEntry{}
	fs := []*model.ContainerInstanceFsEntry{}
	net := []*model.ContainerInstanceNetEntry{}
	for _, input := range op.Inputs {
		switch {
		case nil != input.String:
			{
				stringInput := input.String
				inputValue := ""
				if _, isArgForInput := args[stringInput.Name]; isArgForInput {
					// use provided arg for param
					inputValue = args[stringInput.Name].String
				} else {
					// no provided arg for param; fallback to default
					inputValue = stringInput.Default
				}
				for _, rawCmdEntry := range containerCall.Cmd {
					cmd = append(cmd, interpolater.Interpolate(rawCmdEntry, stringInput.Name, inputValue))
				}
				fmt.Printf("%v\n", cmd)
				for _, envEntry := range containerCall.Env {
					// append bound strings to env
					if envEntry.Bind == stringInput.Name {
						env = append(
							env,
							&model.ContainerInstanceEnvEntry{
								Name:  stringInput.Name,
								Value: inputValue,
							},
						)
						break
					}
				}
			}
		case nil != input.Dir:
			{
				dirInput := input.Dir
				for _, fsEntry := range containerCall.Fs {
					// append bound dirs to fs
					if fsEntry.Bind == dirInput.Name {
						fs = append(
							fs,
							&model.ContainerInstanceFsEntry{
								SrcRef: args[dirInput.Name].Dir,
								Path:   fsEntry.Path,
							},
						)
						break
					}
				}
			}
		case nil != input.File:
			{
				fileInput := input.File
				for _, fsEntry := range containerCall.Fs {
					// append bound files to fs
					if fsEntry.Bind == fileInput.Name {
						fs = append(
							fs,
							&model.ContainerInstanceFsEntry{
								SrcRef: args[fileInput.Name].File,
								Path:   fsEntry.Path,
							},
						)
						break
					}
				}
			}
		case nil != input.NetSocket:
			{
				netSocketInput := input.NetSocket
				netSocketArg := args[netSocketInput.Name].NetSocket
				for _, netEntry := range containerCall.Net {
					// append bound sockets to net
					if netEntry.Bind == netSocketInput.Name {
						net = append(
							net,
							&model.ContainerInstanceNetEntry{
								Host:        netSocketArg.Host,
								Port:        netSocketArg.Port,
								HostAliases: netEntry.HostAliases,
							},
						)
						break
					}
				}
			}
		}
	}

	err = this.containerEngine.StartContainer(
		cmd,
		env,
		fs,
		containerCall.Image,
		net,
		containerCall.WorkDir,
		containerId,
		this.eventBus,
		opGraphId,
	)

	defer func() {

		this.nodeRepo.deleteIfExists(containerId)

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
