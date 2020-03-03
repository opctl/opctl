package hostruntime

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"context"
	"github.com/docker/docker/api/types"
)

//counterfeiter:generate -o internal/fakes/containerInspector.go . containerInspector
type containerInspector interface {
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
}
