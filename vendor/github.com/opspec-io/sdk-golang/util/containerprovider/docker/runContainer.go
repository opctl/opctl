package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"io"
	"strings"
	"sync"
)

type runContainer interface {
	RunContainer(
		req *model.DCGContainerCall,
		eventPublisher pubsub.EventPublisher,
		stdout io.WriteCloser,
		stderr io.WriteCloser,
	) (*int64, error)
}

func newRunContainer(
	dockerClient dockerClientPkg.CommonAPIClient,
) runContainer {
	return _runContainer{
		containerConfigFactory:  newContainerConfigFactory(),
		containerStdErrStreamer: newContainerStdErrStreamer(dockerClient),
		containerStdOutStreamer: newContainerStdOutStreamer(dockerClient),
		dockerClient:            dockerClient,
		hostConfigFactory:       newHostConfigFactory(),
		imagePuller:             newImagePuller(dockerClient),
		portBindingsFactory:     newPortBindingsFactory(),
	}
}

type _runContainer struct {
	containerConfigFactory  containerConfigFactory
	containerStdErrStreamer containerLogStreamer
	containerStdOutStreamer containerLogStreamer
	dockerClient            dockerClientPkg.CommonAPIClient
	hostConfigFactory       hostConfigFactory
	imagePuller             imagePuller
	portBindingsFactory     portBindingsFactory
}

func (cr _runContainer) RunContainer(
	req *model.DCGContainerCall,
	eventPublisher pubsub.EventPublisher,
	stdout io.WriteCloser,
	stderr io.WriteCloser,
) (*int64, error) {
	defer stdout.Close()
	defer stderr.Close()

	portBindings, err := cr.portBindingsFactory.Construct(
		req.Ports,
	)
	if nil != err {
		return nil, err
	}

	containerConfig := cr.containerConfigFactory.Construct(
		req.Cmd,
		req.EnvVars,
		req.Image.Ref,
		portBindings,
		req.WorkDir,
	)

	hostConfig := cr.hostConfigFactory.Construct(
		req.Dirs,
		req.Files,
		req.Sockets,
		portBindings,
	)

	// construct networking config
	networkingConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			dockerNetworkName: {
				Aliases: []string{
					req.Name,
				},
			},
		},
	}

	// always pull latest version of image
	// note: this trades local reproducibility for distributed reproducibility
	pullErr := cr.imagePuller.Pull(
		req.Image,
		req.ContainerId,
		req.RootOpId,
		eventPublisher,
	)
	// don't err yet; image might be cached

	// create container
	containerCreatedResponse, createErr := cr.dockerClient.ContainerCreate(
		context.TODO(),
		containerConfig,
		hostConfig,
		networkingConfig,
		req.ContainerId,
	)
	if nil != createErr {
		if nil == pullErr {
			return nil, createErr
		}
		// if pullErr occurred prior; combine errors
		return nil, errors.New(strings.Join([]string{pullErr.Error(), createErr.Error()}, ", "))
	}

	// start container
	if err := cr.dockerClient.ContainerStart(
		context.TODO(),
		containerCreatedResponse.ID,
		types.ContainerStartOptions{},
	); nil != err {
		return nil, err
	}

	var waitGroup sync.WaitGroup
	errChan := make(chan error, 3)
	waitGroup.Add(2)

	go func() {
		if err := cr.containerStdErrStreamer.Stream(
			req.ContainerId,
			stderr,
		); nil != err {
			errChan <- err
		}
		waitGroup.Done()
	}()

	go func() {
		if err := cr.containerStdOutStreamer.Stream(
			req.ContainerId,
			stdout,
		); nil != err {
			errChan <- err
		}
		waitGroup.Done()
	}()

	var exitCode int64
	waitOkChan, waitErrChan := cr.dockerClient.ContainerWait(
		context.TODO(),
		req.ContainerId,
		container.WaitConditionNotRunning,
	)
	select {
	case waitOk := <-waitOkChan:
		exitCode = waitOk.StatusCode
	case waitErr := <-waitErrChan:
		err = fmt.Errorf("error encountered waiting on container; error was: %v", waitErr.Error())
	}

	// ensure stdout, and stderr all read before returning
	waitGroup.Wait()

	if nil != err && len(errChan) > 0 {
		// non-destructively set err
		err = <-errChan
	}
	return &exitCode, err

}
