package containerengine

import "github.com/opspec-io/sdk-golang/pkg/model"

type StartContainerReq struct {
	Cmd         []string
	Env         []*model.ContainerEnvEntry
	Fs          []*model.ContainerFsEntry
	Image       string
	Net         []*model.ContainerNetEntry
	WorkDir     string
	ContainerId string
	OpGraphId   string
}
