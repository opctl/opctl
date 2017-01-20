package docker

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"golang.org/x/net/context"
	"os"
	"strings"
)

func (this _containerEngine) InspectContainerIfExists(
	containerId string,
) (container *model.DcgContainerCall, err error) {
	dockerContainer, err := this.dockerClient.ContainerInspect(
		context.Background(),
		containerId,
	)
	if nil != err {
		// ignore not found errors
		if client.IsErrContainerNotFound(err) {
			err = nil
		} else {
			return
		}
	}

	// construct dynamic call graph container from docker container
	container = &model.DcgContainerCall{
		Cmd:     dockerContainer.Config.Entrypoint,
		Dirs:    map[string]string{},
		EnvVars: map[string]string{},
		Files:   map[string]string{},
		Image:   dockerContainer.Image,
		Sockets: map[string]string{},
		WorkDir: dockerContainer.Config.WorkingDir,
	}
	// construct envVars
	for _, rawEnvVar := range dockerContainer.Config.Env {
		rawEnvVarParts := strings.SplitN(rawEnvVar, "=", 2)
		container.EnvVars[rawEnvVarParts[0]] = rawEnvVarParts[1]
	}
	// construct dirs, sockets backed by unix sockets, & files
	for _, mount := range dockerContainer.Mounts {
		var fileInfo os.FileInfo
		fileInfo, err = this.vfs.Stat(mount.Source)
		if nil != err {
			err = nil
			fmt.Printf("Mount not available on opctl host and will be ignored. Mount was: %v \n", mount.Source)
			continue
		}
		if fileInfo.IsDir() {
			container.Dirs[mount.Destination] = mount.Source
		} else if (fileInfo.Mode() & (os.ModeSocket | os.ModeNamedPipe)) != 0 {
			container.Sockets[mount.Destination] = mount.Source
		} else {
			container.Files[mount.Destination] = mount.Source
		}
	}

	return
}
