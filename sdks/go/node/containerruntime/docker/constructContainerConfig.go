package docker

import (
	"fmt"
	"sort"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
)

func constructContainerConfig(
	cmd []string,
	envVars map[string]string,
	imageRef string,
	portBindings nat.PortMap,
	workDir string,
) *container.Config {
	containerConfig := &container.Config{
		Image:        imageRef,
		WorkingDir:   workDir,
		Tty:          true,
		ExposedPorts: nat.PortSet{},
	}

	for _, cmd := range cmd {
		containerConfig.Entrypoint = append(containerConfig.Entrypoint, cmd)
	}

	for port := range portBindings {
		containerConfig.ExposedPorts[port] = struct{}{}
	}

	for envVarName, envVarValue := range envVars {
		containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%v=%v", envVarName, envVarValue))
	}

	// propagate proxy env vars from the opctl node's process environment, without
	// overriding values the op explicitly set, so containers can reach the network
	// on hosts whose only egress route is an HTTP/HTTPS forward proxy.
	for envVarName, envVarValue := range containerruntime.ProxyEnvVars(envVars) {
		containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%v=%v", envVarName, envVarValue))
	}
	// sort binds to make order deterministic; useful for testing
	sort.Strings(containerConfig.Env)

	return containerConfig
}
