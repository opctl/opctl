package docker

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/opspec-io/opctl/util/eventbus"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"golang.org/x/net/context"
)

func (this _containerEngine) StartContainer(
	cmd []string,
	env []*model.ContainerInstanceEnvEntry,
	fs []*model.ContainerInstanceFsEntry,
	image string,
	net []*model.ContainerInstanceNetEntry,
	workDir string,
	containerId string,
	eventPublisher eventbus.EventPublisher,
	opGraphId string,
) (err error) {

	// construct container config
	containerConfig := &container.Config{
		Entrypoint: cmd,
		Image:      image,
		WorkingDir: workDir,
		Tty:        true,
	}
	for _, envEntry := range env {
		containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%v=%v", envEntry.Name, envEntry.Value))
	}

	// construct host config
	hostConfig := &container.HostConfig{
		// support docker in docker
		// @TODO: reconsider; can we avoid this?
		// see for similar discussion: https://github.com/kubernetes/kubernetes/issues/391
		Privileged: true,
	}
	for _, fsEntry := range fs {
		hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%v:%v", fsEntry.SrcRef, fsEntry.Path))
	}
	for _, netEntry := range net {
		for _, hostAlias := range netEntry.HostAliases {
			hostConfig.ExtraHosts = append(hostConfig.ExtraHosts, fmt.Sprintf("%v:%v", hostAlias, netEntry.Host))
		}
	}

	// create container
	var containerCreatedResponse container.ContainerCreateCreatedBody
	containerCreatedResponse, err = this.dockerEngine.ContainerCreate(
		context.Background(),
		containerConfig,
		hostConfig,
		&network.NetworkingConfig{},
		containerId,
	)

	if nil != err {
		//if image not found try to pull it
		if client.IsErrImageNotFound(err) {
			err = nil
			fmt.Printf("Unable to find image '%s' locally\n", image)

			err = this.pullImage(image, containerId, opGraphId, eventPublisher)
			if nil != err {
				return
			}

			// Retry
			containerCreatedResponse, err = this.dockerEngine.ContainerCreate(
				context.Background(),
				containerConfig,
				hostConfig,
				&network.NetworkingConfig{},
				containerId,
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

	err = this.stdErrLogger(eventPublisher, containerId, opGraphId)
	if nil != err {
		return
	}
	err = this.stdOutLogger(eventPublisher, containerId, opGraphId)
	if nil != err {
		return
	}

	// wait for exit
	exitCode, exitCodeReadError := this.dockerEngine.ContainerWait(
		context.Background(),
		containerId,
	)
	if nil != exitCodeReadError {
		err = errors.New(fmt.Sprintf("unable to read container exit code. Error was: %v", exitCodeReadError))
	} else if 0 != exitCode {
		err = errors.New(fmt.Sprintf("nonzero container exit code. Exit code was: %v", exitCode))
	}

	return

}
