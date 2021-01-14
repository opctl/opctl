package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dockerClientPkg "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

//counterfeiter:generate -o internal/fakes/hostConfigFactory.go . hostConfigFactory
type hostConfigFactory interface {
	Construct(
		containerCallDirs map[string]string,
		containerCallFiles map[string]string,
		containerCallSockets map[string]string,
		portBindings nat.PortMap,
	) *container.HostConfig
}

func newHostConfigFactory(
	ctx context.Context,
	dockerClient dockerClientPkg.CommonAPIClient,
) (hostConfigFactory, error) {
	fspc, err := newFSPathConverter(ctx, dockerClient)
	if err != nil {
		return _hostConfigFactory{}, err
	}

	hcf := _hostConfigFactory{
		fsPathConverter: fspc,
	}
	return hcf, nil
}

type _hostConfigFactory struct {
	fsPathConverter fsPathConverter
}

func (hcf _hostConfigFactory) Construct(
	containerCallDirs map[string]string,
	containerCallFiles map[string]string,
	containerCallSockets map[string]string,
	portBindings nat.PortMap,
) *container.HostConfig {
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		// support docker in docker
		// @TODO: reconsider; can we avoid this?
		// see for similar discussion: https://github.com/kubernetes/kubernetes/issues/391
		Privileged: true,
	}
	for containerFilePath, hostFilePath := range containerCallFiles {
		hostConfig.Mounts = append(
			hostConfig.Mounts,
			mount.Mount{
				Type:        mount.TypeBind,
				Source:      hcf.fsPathConverter.LocalToEngine(hostFilePath),
				Target:      containerFilePath,
				Consistency: mount.ConsistencyCached,
			},
		)
	}
	for containerDirPath, hostDirPath := range containerCallDirs {
		hostConfig.Mounts = append(
			hostConfig.Mounts,
			mount.Mount{
				Type:        mount.TypeBind,
				Source:      hcf.fsPathConverter.LocalToEngine(hostDirPath),
				Target:      containerDirPath,
				Consistency: mount.ConsistencyCached,
			},
		)
	}
	for containerSocketAddress, hostSocketAddress := range containerCallSockets {
		const unixSocketAddressDiscriminationChars = `/\`
		// note: this mechanism for determining the type of socket is naive; higher level of sophistication may be required
		if strings.ContainsAny(hostSocketAddress, unixSocketAddressDiscriminationChars) {
			hostConfig.Mounts = append(
				hostConfig.Mounts,
				mount.Mount{
					Type:   mount.TypeBind,
					Source: hcf.fsPathConverter.LocalToEngine(hostSocketAddress),
					Target: containerSocketAddress,
				},
			)
		}
	}

	return hostConfig
}
