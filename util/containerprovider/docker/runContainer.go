package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/pubsub"
	"golang.org/x/net/context"
	"io"
	"sort"
	"strings"
	"sync"
)

func (ctp _containerProvider) RunContainer(
	req *model.DCGContainerCall,
	eventPublisher pubsub.EventPublisher,
	stdout io.WriteCloser,
	stderr io.WriteCloser,
) (*int64, error) {
	defer stdout.Close()
	defer stderr.Close()

	// construct port bindings
	portBindings := nat.PortMap{}
	for containerPort, hostPort := range req.Ports {
		portMappings, err := nat.ParsePortSpec(fmt.Sprintf("%v:%v", hostPort, containerPort))
		if nil != err {
			return nil, err
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
		// @TODO: reconsider; can we avoid ctp?
		// see for similar discussion: https://github.com/kubernetes/kubernetes/issues/391
		Privileged: true,
	}
	for containerFilePath, hostFilePath := range req.Files {
		hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%v:%v", ctp.enginePath(hostFilePath), containerFilePath))
	}
	for containerDirPath, hostDirPath := range req.Dirs {
		hostConfig.Binds = append(hostConfig.Binds, fmt.Sprintf("%v:%v", ctp.enginePath(hostDirPath), containerDirPath))
	}
	for containerSocketAddress, hostSocketAddress := range req.Sockets {
		const unixSocketAddressDiscriminationChars = `/\`
		// note: ctp mechanism for determining the type of socket is naive; higher level of sophistication may be required
		if strings.ContainsAny(hostSocketAddress, unixSocketAddressDiscriminationChars) {
			hostConfig.Binds = append(
				hostConfig.Binds,
				fmt.Sprintf("%v:%v", ctp.enginePath(hostSocketAddress), containerSocketAddress),
			)
		}
	}
	// sort binds to make order deterministic; useful for testing
	sort.Strings(hostConfig.Binds)

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

	// create container
	containerCreatedResponse, err := ctp.dockerClient.ContainerCreate(
		context.TODO(),
		containerConfig,
		hostConfig,
		networkingConfig,
		req.ContainerId,
	)

	if nil != err {
		if !client.IsErrNotFound(err) {
			return nil, err
		}

		// pull image
		err = nil
		fmt.Printf("unable to find image '%s' locally\n", req.Image.Ref)

		err = ctp.pullImage(req.Image, req.ContainerId, req.RootOpId, eventPublisher)
		if nil != err {
			return nil, err
		}

		// retry create
		containerCreatedResponse, err = ctp.dockerClient.ContainerCreate(
			context.TODO(),
			containerConfig,
			hostConfig,
			networkingConfig,
			req.ContainerId,
		)
		if nil != err {
			return nil, err
		}
	}

	// start container
	if err := ctp.dockerClient.ContainerStart(
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
		if err := ctp.streamContainerStdErr(
			req.ContainerId,
			stderr,
		); nil != err {
			errChan <- err
		}
		waitGroup.Done()
	}()

	go func() {
		if err := ctp.streamContainerStdOut(
			req.ContainerId,
			stdout,
		); nil != err {
			errChan <- err
		}
		waitGroup.Done()
	}()

	var exitCode int64
	waitOkChan, waitErrChan := ctp.dockerClient.ContainerWait(
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
