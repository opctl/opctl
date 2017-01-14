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
	rawContainer, err := this.dockerEngine.ContainerInspect(
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

	// construct container from rawContainer
	container = &model.DcgContainerCall{
		Cmd:     rawContainer.Config.Entrypoint,
		EnvVars: map[string]string{},
		Dirs:    map[string]string{},
		Files:   map[string]string{},
		WorkDir: rawContainer.Config.WorkingDir,
		Image:   rawContainer.Image,
	}
	// construct envVars
	for _, rawEnvVar := range rawContainer.Config.Env {
		rawEnvVarParts := strings.SplitN(rawEnvVar, "=", 2)
		container.EnvVars[rawEnvVarParts[0]] = rawEnvVarParts[1]
	}
	// construct files & dirs
	for _, mount := range rawContainer.Mounts {
		var fileInfo os.FileInfo
		fileInfo, err = this.fs.Stat(mount.Source)
		if nil != err {
			err = nil
			fmt.Printf("Mount not available on opctl host and will be ignored. Mount was: %v \n", mount.Source)
			continue
		}
		if fileInfo.IsDir() {
			container.Dirs[mount.Destination] = mount.Source
		} else {
			container.Files[mount.Destination] = mount.Source
		}
	}
	// @todo: construct sockets

	return
}
