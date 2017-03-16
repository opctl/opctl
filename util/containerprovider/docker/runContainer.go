package docker

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/opspec-io/opctl/util/pubsub"
	"github.com/opspec-io/sdk-golang/model"
	"golang.org/x/net/context"
	"sort"
	"strings"
)

func (this _containerProvider) RunContainer(
	req *model.DCGContainerCall,
	eventPublisher pubsub.EventPublisher,
) (err error) {

	// construct port bindings
	portBindings := nat.PortMap{}
	for containerPort, hostPort := range req.Ports {
		portMappings, err := nat.ParsePortSpec(fmt.Sprintf("%v:%v", hostPort, containerPort))
		if nil != err {
			return err
		}
		for _, portMapping := range portMappings {
			if _, ok := portBindings[portMapping.Port]; ok {
				portBindings[portMapping.Port] = append(portBindings[portMapping.Port], portMapping.Binding)
			} else {
				portBindings[portMapping.Port] = []nat.PortBinding{portMapping.Binding}
			}
		}
	}

	// construct container config
	containerConfig := &container.Config{
		Entrypoint:   req.Cmd,
		Image:        req.Image.Ref,
		WorkingDir:   req.WorkDir,
		Tty:          true,
		ExposedPorts: nat.PortSet{},
	}
	for port := range portBindings {
		containerConfig.ExposedPorts[port] = struct{}{}
	}
	for envVarName, envVarValue := range req.EnvVars {
		containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%v=%v", envVarName, envVarValue))
	}
	// sort binds to make order deterministic; useful for testing
	sort.Strings(containerConfig.Env)

	// construct host config
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		// support docker in docker
		// @TODO: reconsider; can we avoid this?
		// see for similar discussion: https://github.com/kubernetes/kubernetes/issues/391
		Privileged: true,
	}
	for containerFilePath, hostFilePath := range req.Files {
		hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%v:%v", this.enginePath(hostFilePath), containerFilePath))
	}
	for containerDirPath, hostDirPath := range req.Dirs {
		hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%v:%v", this.enginePath(hostDirPath), containerDirPath))
	}
	for containerSocketAddress, hostSocketAddress := range req.Sockets {
		const unixSocketAddressDiscriminationChars = `/\`
		// note: this mechanism for determining the type of socket is naive; higher level of sophistication may be required
		if strings.ContainsAny(hostSocketAddress, unixSocketAddressDiscriminationChars) {
			hostConfig.Binds = append(
				hostConfig.Binds,
				fmt.Sprintf("%v:%v", this.enginePath(hostSocketAddress), containerSocketAddress),
			)
		}
	}
	// sort binds to make order deterministic; useful for testing
	sort.Strings(hostConfig.Binds)

	// create container
	containerCreatedResponse, err := this.dockerClient.ContainerCreate(
		context.Background(),
		containerConfig,
		hostConfig,
		&network.NetworkingConfig{
			EndpointsConfig: map[string]*network.EndpointSettings{
				"opctl": {
					Aliases: []string{
						req.Name,
					},
				},
			},
		},
		req.ContainerId,
	)

	if nil != err {
		//if image not found try to pull it
		if client.IsErrImageNotFound(err) {
			err = nil
			fmt.Printf("Unable to find image '%s' locally\n", req.Image.Ref)

			err = this.pullImage(req.Image, req.ContainerId, req.RootOpId, eventPublisher)
			if nil != err {
				return
			}

			// Retry
			containerCreatedResponse, err = this.dockerClient.ContainerCreate(
				context.Background(),
				containerConfig,
				hostConfig,
				&network.NetworkingConfig{},
				req.ContainerId,
			)
			if nil != err {
				return
			}
		} else {
			return
		}
	}

	// start container
	err = this.dockerClient.ContainerStart(
		context.Background(),
		containerCreatedResponse.ID,
		types.ContainerStartOptions{},
	)
	if nil != err {
		return
	}

	err = this.stdErrLogger(eventPublisher, req.ContainerId, req.RootOpId)
	if nil != err {
		return
	}
	err = this.stdOutLogger(eventPublisher, req.ContainerId, req.RootOpId)
	if nil != err {
		return
	}

	// wait for exit
	exitCode, waitError := this.dockerClient.ContainerWait(
		context.Background(),
		req.ContainerId,
	)
	if nil != waitError {
		err = errors.New(fmt.Sprintf("unable to read container exit code. Error was: %v", waitError))
	} else if 0 != exitCode {
		err = errors.New(fmt.Sprintf("nonzero container exit code. Exit code was: %v", exitCode))
	}

	return

}
