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
		req *model.ContainerCall,
		rootCallID string,
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
		ensureNetworkExistser:   newEnsureNetworkExistser(dockerClient),
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
	ensureNetworkExistser   ensureNetworkExistser
	hostConfigFactory       hostConfigFactory
	imagePuller             imagePuller
	imagePusher             imagePusher
	portBindingsFactory     portBindingsFactory
}

func (cr _runContainer) RunContainer(
	ctx context.Context,
	req *model.ContainerCall,
	rootCallID string,
	eventPublisher pubsub.EventPublisher,
	stdout io.WriteCloser,
	stderr io.WriteCloser,
) (*int64, error) {
	defer stdout.Close()
	defer stderr.Close()

	// ensure user defined network exists to allow inter container resolution via name
	// @TODO: remove when socket outputs supported
	if err := cr.ensureNetworkExistser.EnsureNetworkExists(
		ctx,
		dockerNetworkName,
	); nil != err {
		return nil, err
	}

	// for docker, we prefix name with opctl_ in order to allow external tools to know it's an opctl managed container
	// do not change this prefix as it might break external consumers
	containerName := fmt.Sprintf("opctl_%s", req.ContainerID)
	defer func() {
		// ensure container always cleaned up
		cr.dockerClient.ContainerRemove(
			context.Background(),
			containerName,
			types.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			},
		)
	}()

	var imageErr error
	if nil != req.Image.Src {
		imageRef := fmt.Sprintf("%s:latest", req.ContainerID)
		req.Image.Ref = &imageRef

		imageErr = cr.imagePusher.Push(
			ctx,
			imageRef,
			req.Image.Src,
		)
	} else {
		// always pull latest version of image
		// note: this trades local reproducibility for distributed reproducibility
		imageErr = cr.imagePuller.Pull(
			ctx,
			req.ContainerID,
			req.Image.PullCreds,
			*req.Image.Ref,
			rootCallID,
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
		containerName,
	)
	if nil != createErr {
		select {
		case <-ctx.Done():
			// we got killed;
			return nil, nil
		default:
			if nil == imageErr {
				return nil, createErr
			}
			// if imageErr occurred prior; combine errors
			return nil, errors.New(strings.Join([]string{imageErr.Error(), createErr.Error()}, ", "))
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
			containerName,
			stderr,
		); nil != err {
			errChan <- err
		}
		waitGroup.Done()
	}()

	go func() {
		if err := cr.containerStdOutStreamer.Stream(
			ctx,
			containerName,
			stdout,
		); nil != err {
			errChan <- err
		}
		waitGroup.Done()
	}()

	var exitCode int64
	waitOkChan, waitErrChan := cr.dockerClient.ContainerWait(
		ctx,
		containerName,
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
