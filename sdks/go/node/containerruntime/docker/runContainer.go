package docker

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/dns"
	"github.com/opctl/opctl/sdks/go/node/pubsub"
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
	ctx context.Context,
	dockerClient dockerClientPkg.CommonAPIClient,
) (runContainer, error) {
	rc := _runContainer{
		containerStdErrStreamer: newContainerStdErrStreamer(dockerClient),
		containerStdOutStreamer: newContainerStdOutStreamer(dockerClient),
		dockerClient:            dockerClient,
	}
	return rc, nil
}

type _runContainer struct {
	containerStdErrStreamer containerLogStreamer
	containerStdOutStreamer containerLogStreamer
	dockerClient            dockerClientPkg.CommonAPIClient
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
	var containerIPAddress string

	// ensure user defined network exists to allow inter container resolution via name
	// @TODO: remove when socket outputs supported
	if err := ensureNetworkExists(
		ctx,
		cr.dockerClient,
		req.Image.PullCreds,
		networkName,
	); err != nil {
		return nil, err
	}

	// for docker, we prefix name with opctl_ in order to allow external tools to know it's an opctl managed container
	// do not change this prefix as it might break external consumers
	dockerContainerName := getContainerName(req.ContainerID)
	defer func() {
		// ensure container always cleaned up: gracefully stop then delete it
		// @TODO: consolidate logic with DeleteContainerIfExists
		newCtx := context.Background() // always use a fresh context, to clean up after cancellation
		stopTimeout := 3
		cr.dockerClient.ContainerStop(
			newCtx,
			dockerContainerName,
			container.StopOptions{
				Timeout: &stopTimeout,
			},
		)
		cr.dockerClient.ContainerRemove(
			newCtx,
			dockerContainerName,
			container.RemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			},
		)

		if req.Name != nil {
			dns.UnregisterName(
				*req.Name,
				containerIPAddress,
			)
		}
	}()

	var imageErr error
	// if req.Image.Src != nil {
	// 	imageRef := fmt.Sprintf("%s:latest", req.ContainerID)
	// 	req.Image.Ref = &imageRef

	// 	imageErr = pushImage(
	// 		ctx,
	// 		imageRef,
	// 		req.Image.Src,
	// 	)
	// } else {
	imageErr = pullImage(
		ctx,
		req,
		cr.dockerClient,
		rootCallID,
		eventPublisher,
	)
	// don't err yet; image might be cached. We allow this to support offline use
	// }

	portBindings, err := constructPortBindings(
		req.Ports.Values,
	)
	if err != nil {
		return nil, err
	}

	// construct networking config
	networkingConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			networkName: {},
		},
	}
	if req.Name != nil {
		networkingConfig.EndpointsConfig[networkName].Aliases = []string{
			*req.Name,
		}
	}

	isGpuSupported, err := isGpuSupported(ctx, cr.dockerClient, req.Image.PullCreds)
	if nil != err {
		return nil, err
	}

	// create container
	_, createErr := cr.dockerClient.ContainerCreate(
		ctx,
		constructContainerConfig(
			req.Cmd,
			req.EnvVars.Values,
			*req.Image.Ref,
			portBindings,
			req.WorkDir,
		),
		constructHostConfig(
			req.Dirs.Values,
			req.Files.Values,
			req.Sockets.Values,
			portBindings,
			isGpuSupported,
		),
		networkingConfig,
		// platform requires API v1.41 so set to nil to avoid version errors
		nil,
		dockerContainerName,
	)
	if createErr != nil {
		select {
		case <-ctx.Done():
			// we got killed;
			return nil, nil
		default:
			if imageErr == nil {
				return nil, createErr
			}
			// if imageErr occurred prior; combine errors
			return nil, errors.New(strings.Join([]string{imageErr.Error(), createErr.Error()}, ", "))
		}
	}

	// start container
	if err := cr.dockerClient.ContainerStart(
		ctx,
		dockerContainerName,
		container.StartOptions{},
	); err != nil {
		return nil, err
	}

	containerJSON, err := cr.dockerClient.ContainerInspect(ctx, dockerContainerName)
	if err != nil {
		return nil, err
	}

	if endpointSettings, ok := containerJSON.NetworkSettings.Networks[networkName]; ok && req.Name != nil {
		err = dns.RegisterName(
			*req.Name,
			endpointSettings.IPAddress,
		)
		if err != nil {
			return nil, err
		}
	}

	var waitGroup sync.WaitGroup
	errChan := make(chan error, 3)
	waitGroup.Add(2)

	go func() {
		if err := cr.containerStdErrStreamer.Stream(
			ctx,
			dockerContainerName,
			stderr,
		); err != nil {
			errChan <- err
		}
		waitGroup.Done()
	}()

	go func() {
		if err := cr.containerStdOutStreamer.Stream(
			ctx,
			dockerContainerName,
			stdout,
		); err != nil {
			errChan <- err
		}
		waitGroup.Done()
	}()

	var exitCode int64
	waitOkChan, waitErrChan := cr.dockerClient.ContainerWait(
		ctx,
		dockerContainerName,
		container.WaitConditionNotRunning,
	)
	select {
	case waitOk := <-waitOkChan:
		exitCode = waitOk.StatusCode
	case waitErr := <-waitErrChan:
		err = fmt.Errorf("error waiting on container: %w", waitErr)
	}

	// ensure stdout, and stderr all read before returning
	waitGroup.Wait()

	if err != nil && len(errChan) > 0 {
		// non-destructively set err
		err = <-errChan
	}
	return &exitCode, err

}
