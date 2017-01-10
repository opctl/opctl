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
		Volumes:    map[string]struct{}{},
	}
	for _, envEntry := range req.Env {
		containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%v=%v", envEntry.Name, envEntry.Value))
	}

	// construct host config
	hostConfig := &container.HostConfig{
		// support docker in docker
		// @TODO: reconsider; can we avoid this?
		// see for similar discussion: https://github.com/kubernetes/kubernetes/issues/391
		Privileged: true,
	}
	for _, fsEntry := range req.Fs {
		if "" != fsEntry.SrcRef {
			// bind mount
			hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%v:%v", fsEntry.SrcRef, fsEntry.Path))
		} else {
			// anonymous volume
			containerConfig.Volumes[fsEntry.Path] = struct{}{}
		}
	}
	for _, netEntry := range req.Net {
		for _, hostAlias := range netEntry.HostAliases {
			hostConfig.ExtraHosts = append(hostConfig.ExtraHosts, fmt.Sprintf("%v:%v", hostAlias, netEntry.Host))
		}
	}

	fmt.Printf("startContainer: hostConfig\n%#v\n", hostConfig)

	// create container
	var containerCreatedResponse container.ContainerCreateCreatedBody
	containerCreatedResponse, err = this.dockerEngine.ContainerCreate(
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
			containerCreatedResponse, err = this.dockerEngine.ContainerCreate(
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
	err = this.dockerEngine.ContainerStart(
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
	exitCode, waitError := this.dockerEngine.ContainerWait(
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
