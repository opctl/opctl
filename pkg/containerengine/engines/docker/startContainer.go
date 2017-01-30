package docker

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/opspec-io/opctl/pkg/containerengine"
	"github.com/opspec-io/opctl/util/eventbus"
	"golang.org/x/net/context"
	"sort"
	"strings"
)

func (this _containerEngine) StartContainer(
	req *containerengine.StartContainerReq,
	eventPublisher eventbus.EventPublisher,
) (err error) {

	// construct container config
	containerConfig := &container.Config{
		Entrypoint: req.Cmd,
		Image:      req.Image,
		WorkingDir: req.WorkDir,
		Tty:        true,
	}
	for envVarName, envVarValue := range req.Env {
		containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%v=%v", envVarName, envVarValue))
	}
	// sort binds to make order deterministic; useful for testing
	sort.Strings(containerConfig.Env)

	// construct host config
	hostConfig := &container.HostConfig{
		// support docker in docker
		// @TODO: reconsider; can we avoid this?
		// see for similar discussion: https://github.com/kubernetes/kubernetes/issues/391
		Privileged: true,
	}
	for containerFilePath, hostFilePath := range req.Files {
		hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%v:%v", hostFilePath, containerFilePath))
	}
	for containerDirPath, hostDirPath := range req.Dirs {
		hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%v:%v", hostDirPath, containerDirPath))
	}
	for containerSocketAddress, hostSocketAddress := range req.Sockets {
		const unixSocketAddressDiscriminationChars = `/\`
		switch {
		// note: this mechanism for determining the type of socket is naive; higher level of sophistication may be required
		case strings.ContainsAny(hostSocketAddress, unixSocketAddressDiscriminationChars):
			hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%v:%v", hostSocketAddress, containerSocketAddress))
		default:
			// @TODO: handle network sockets
		}
	}
	// sort binds to make order deterministic; useful for testing
	sort.Strings(hostConfig.Binds)

	// create container
	containerCreatedResponse, err := this.dockerClient.ContainerCreate(
		context.Background(),
		containerConfig,
		hostConfig,
		&network.NetworkingConfig{},
		req.ContainerId,
	)

	if nil != err {
		//if image not found try to pull it
		if client.IsErrImageNotFound(err) {
			err = nil
			fmt.Printf("Unable to find image '%s' locally\n", req.Image)

			err = this.pullImage(req.Image, req.ContainerId, req.OpGraphId, eventPublisher)
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

	err = this.stdErrLogger(eventPublisher, req.ContainerId, req.OpGraphId)
	if nil != err {
		return
	}
	err = this.stdOutLogger(eventPublisher, req.ContainerId, req.OpGraphId)
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
