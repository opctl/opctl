package docker

import (
	"fmt"
	"sort"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

//counterfeiter:generate -o internal/fakes/containerConfigFactory.go . containerConfigFactory
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
	// sort binds to make order deterministic; useful for testing
	sort.Strings(containerConfig.Env)

	return containerConfig
}
