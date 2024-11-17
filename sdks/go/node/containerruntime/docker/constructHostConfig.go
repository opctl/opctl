package docker

import (
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

func constructHostConfig(
	containerCallDirs map[string]string,
	containerCallFiles map[string]string,
	containerCallSockets map[string]string,
	portBindings nat.PortMap,
	isGpuSupported bool,
) *container.HostConfig {
	hostConfig := &container.HostConfig{
		PortBindings: portBindings,
		// support docker in docker
		// @TODO: reconsider; can we avoid this?
		// see for similar discussion: https://github.com/kubernetes/kubernetes/issues/391
		Privileged: true,
	}
	if isGpuSupported {
		hostConfig.Resources = container.Resources{
			DeviceRequests: []container.DeviceRequest{
				{
					Capabilities: [][]string{{"gpu"}},
					Count:        -1,
				},
			},
		}
	}

	for containerFilePath, hostFilePath := range containerCallFiles {
		hostConfig.Mounts = append(
			hostConfig.Mounts,
			mount.Mount{
				Type:        mount.TypeBind,
				Source:      hostFilePath,
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
				Source:      hostDirPath,
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
					Source: hostSocketAddress,
					Target: containerSocketAddress,
				},
			)
		}
	}

	return hostConfig
}
