package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/pubsub"
)

type runContainer interface {
	RunContainer(
		ctx context.Context,
		req *model.DCGContainerCall,
		eventPublisher pubsub.EventPublisher,
		stdout io.WriteCloser,
		stderr io.WriteCloser,
	) (*int64, error)
}

func newRunContainer(
	dockerClient dockerClientPkg.CommonAPIClient,
) (runContainer, error) {
	hcf, err := newHostConfigFactory(dockerClient)
	if err != nil {
		return _runContainer{}, err
	}

	rc := _runContainer{
		containerConfigFactory:  newContainerConfigFactory(),
		containerStdErrStreamer: newContainerStdErrStreamer(dockerClient),
		containerStdOutStreamer: newContainerStdOutStreamer(dockerClient),
		dockerClient:            dockerClient,
		hostConfigFactory:       hcf,
		imagePuller:             newImagePuller(dockerClient),
		imagePusher:             newImagePusher(),
		portBindingsFactory:     newPortBindingsFactory(),
	}
	return rc, nil
}

type _runContainer struct {
	containerConfigFactory  containerConfigFactory
	containerStdErrStreamer containerLogStreamer
	containerStdOutStreamer containerLogStreamer
	dockerClient            dockerClientPkg.CommonAPIClient
	hostConfigFactory       hostConfigFactory
	imagePuller             imagePuller
	imagePusher             imagePusher
	portBindingsFactory     portBindingsFactory
}

func (cr _runContainer) RunContainer(
	ctx context.Context,
	req *model.DCGContainerCall,
	eventPublisher pubsub.EventPublisher,
	stdout io.WriteCloser,
	stderr io.WriteCloser,
) (*int64, error) {
	defer stdout.Close()
	defer stderr.Close()
	defer func() {
		// ensure container always cleaned up
		cr.dockerClient.ContainerRemove(
			context.Background(),
			req.ContainerID,
			types.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			},
		)
	}()

	var pullErr error
	if nil != req.Image.Src {
		imageRef := fmt.Sprintf("%s:latest", req.ContainerID)
		req.Image.Ref = &imageRef

		pullErr = cr.imagePusher.Push(
			ctx,
			req.ContainerID,
			imageRef,
			req.Image.Src,
			req.RootOpID,
			eventPublisher,
		)
	} else {
		// always pull latest version of image
		// note: this trades local reproducibility for distributed reproducibility
		pullErr = cr.imagePuller.Pull(
			ctx,
			req.ContainerID,
			req.Image.PullCreds,
			*req.Image.Ref,
			req.RootOpID,
			eventPublisher,
		)
		// don't err yet; image might be cached. We allow this to support offline use
	}

	portBindings, err := cr.portBindingsFactory.Construct(
		req.Ports,
	)
	if nil != err {
		return nil, err
	}

	containerConfig := cr.containerConfigFactory.Construct(
		req.Cmd,
		req.EnvVars,
		*req.Image.Ref,
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
			dockerNetworkName: {},
		},
	}
	if nil != req.Name {
		networkingConfig.EndpointsConfig[dockerNetworkName].Aliases = []string{
			*req.Name,
		}
	}

	// create container
	containerCreatedResponse, createErr := cr.dockerClient.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		networkingConfig,
		req.ContainerID,
	)
	if nil != createErr {
		select {
		case <-ctx.Done():
			// we got killed;
			return nil, nil
		default:
			if nil == pullErr {
				return nil, createErr
			}
			// if pullErr occurred prior; combine errors
			return nil, errors.New(strings.Join([]string{pullErr.Error(), createErr.Error()}, ", "))
		}
	}

	// start container
	if err := cr.dockerClient.ContainerStart(
		ctx,
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
			ctx,
			req.ContainerID,
			stderr,
		); nil != err {
			errChan <- err
		}
		waitGroup.Done()
	}()

	go func() {
		if err := cr.containerStdOutStreamer.Stream(
			ctx,
			req.ContainerID,
			stdout,
		); nil != err {
			errChan <- err
		}
		waitGroup.Done()
	}()

	var exitCode int64
	waitOkChan, waitErrChan := cr.dockerClient.ContainerWait(
		ctx,
		req.ContainerID,
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
