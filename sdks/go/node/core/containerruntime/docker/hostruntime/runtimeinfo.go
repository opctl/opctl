package hostruntime

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"time"
)

// RuntimeInfo provides relation between opctl and docker engine host
type RuntimeInfo struct {
	// InAContainer indicates if opctl is running on and spinning up new containers on same docker engine host
	InAContainer bool
	// DockerID is ID of opctl's container
	DockerID string
	// HostPathMap provides remapping of paths inside opctl container to paths on docker engine host
	HostPathMap HostPathMap
}

func New(cli containerInspector) (RuntimeInfo, error) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	return newContainerRuntimeInfo(ctx, cli, defaultContainerUtils)
}

func newContainerRuntimeInfo(ctx context.Context, cli containerInspector, cu containerUtils) (RuntimeInfo, error) {
	cri := RuntimeInfo{
		HostPathMap: make(map[string]string),
	}

	if cu.inAContainer() {
		dockerID, err := cu.getDockerID()
		if err != nil {
			return cri, errors.Wrap(err, "can't get container's docker id")
		}

		cri.DockerID = dockerID

		info, exists, err := inspect(ctx, cli, dockerID)
		if err != nil {
			return cri, errors.Wrap(err, "can't get inspect current container")
		}

		if exists {
			cri.InAContainer = true
			cri.HostPathMap = newHostPathMap(info.HostConfig.Binds)
		}
	}

	return cri, nil
}

func inspect(ctx context.Context, cli containerInspector, dockerID string) (types.ContainerJSON, bool, error) {
	info, err := cli.ContainerInspect(ctx, dockerID)
	if err != nil && client.IsErrNotFound(err) {
		return info, false, nil
	}
	if err != nil {
		return info, false, err
	}

	return info, true, nil
}
