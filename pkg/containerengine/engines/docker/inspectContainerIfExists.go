package docker

import (
	"github.com/docker/docker/client"
	"github.com/opspec-io/sdk-golang/pkg/model"
	"golang.org/x/net/context"
	"strings"
)

func (this _containerEngine) InspectContainerIfExists(
	containerId string,
) (container *model.Container, err error) {
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
	container = &model.Container{
		Cmd:     rawContainer.Config.Entrypoint,
		WorkDir: rawContainer.Config.WorkingDir,
		Image:   rawContainer.Image,
	}
	// construct env
	for _, rawEnvEntry := range rawContainer.Config.Env {
		rawEnvEntryParts := strings.SplitN(rawEnvEntry, "=", 2)
		container.Env = append(
			container.Env,
			&model.ContainerEnvEntry{
				Name:  rawEnvEntryParts[0],
				Value: rawEnvEntryParts[1],
			},
		)
	}
	// construct fs
	for _, rawFsEntry := range rawContainer.Mounts {
		container.Fs = append(
			container.Fs,
			&model.ContainerFsEntry{
				SrcRef: rawFsEntry.Source,
				Path:   rawFsEntry.Destination,
			},
		)
	}
	// @todo: construct net

	return
}
