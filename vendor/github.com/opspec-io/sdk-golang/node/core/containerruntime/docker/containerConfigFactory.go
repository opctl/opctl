package docker

//go:generate counterfeiter -o ./fakeContainerConfigFactory.go --fake-name fakeContainerConfigFactory ./ containerConfigFactory

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"sort"
)

type containerConfigFactory interface {
	Construct(
		cmd []string,
		envVars map[string]string,
		imageRef string,
		portBindings nat.PortMap,
		workDir string,
	) *container.Config
}

func newContainerConfigFactory() containerConfigFactory {
	return _containerConfigFactory{}
}

type _containerConfigFactory struct{}

func (ccf _containerConfigFactory) Construct(
	cmd []string,
	envVars map[string]string,
	imageRef string,
	portBindings nat.PortMap,
	workDir string,
) *container.Config {
	containerConfig := &container.Config{
		Entrypoint:   cmd,
		Image:        imageRef,
		WorkingDir:   workDir,
		Tty:          true,
		ExposedPorts: nat.PortSet{},
	}
	for port := range portBindings {
		containerConfig.ExposedPorts[port] = struct{}{}
	}
	for envVarName, envVarValue := range envVars {
		containerConfig.Env = append(containerConfig.Env, fmt.Sprintf("%v=%v", envVarName, envVarValue))
	}
	// sort binds to make order deterministic; useful for testing
	sort.Strings(containerConfig.Env)

	return containerConfig
}
