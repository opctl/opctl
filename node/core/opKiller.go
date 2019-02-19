package core

//go:generate counterfeiter -o ./fakeOpKiller.go --fake-name fakeOpKiller ./ opKiller

import (
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/node/core/containerruntime"
	"sync"
)

type opKiller interface {
	Kill(req model.KillOpReq)
}

func newOpKiller(
	dcgNodeRepo dcgNodeRepo,
	containerRuntime containerruntime.ContainerRuntime,
) opKiller {
	return _opKiller{
		dcgNodeRepo:      dcgNodeRepo,
		containerRuntime: containerRuntime,
	}
}

type _opKiller struct {
	dcgNodeRepo      dcgNodeRepo
	containerRuntime containerruntime.ContainerRuntime
}

func (this _opKiller) Kill(
	req model.KillOpReq,
) {
	this.dcgNodeRepo.DeleteIfExists(req.OpID)

	var waitGroup sync.WaitGroup

	for _, childNode := range this.dcgNodeRepo.ListWithRootOpID(req.OpID) {
		waitGroup.Add(1)
		go func(childNode *dcgNodeDescriptor) {
			this.dcgNodeRepo.DeleteIfExists(childNode.Id)

			if nil != childNode.Container {
				this.containerRuntime.DeleteContainerIfExists(
					childNode.Id,
				)
			}
			defer waitGroup.Done()
		}(childNode)
	}

	waitGroup.Wait()

}
