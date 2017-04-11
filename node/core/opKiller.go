package core

//go:generate counterfeiter -o ./fakeOpKiller.go --fake-name fakeOpKiller ./ opKiller

import (
	"github.com/opctl/opctl/util/containerprovider"
	"github.com/opspec-io/sdk-golang/model"
	"sync"
)

type opKiller interface {
	Kill(req model.KillOpReq)
}

func newOpKiller(
	dcgNodeRepo dcgNodeRepo,
	containerProvider containerprovider.ContainerProvider,
) opKiller {
	return _opKiller{
		dcgNodeRepo:       dcgNodeRepo,
		containerProvider: containerProvider,
	}
}

type _opKiller struct {
	dcgNodeRepo       dcgNodeRepo
	containerProvider containerprovider.ContainerProvider
}

func (this _opKiller) Kill(
	req model.KillOpReq,
) {
	this.dcgNodeRepo.DeleteIfExists(req.OpId)

	var waitGroup sync.WaitGroup

	for _, childNode := range this.dcgNodeRepo.ListWithRootOpId(req.OpId) {
		waitGroup.Add(1)
		go func(childNode *dcgNodeDescriptor) {
			this.dcgNodeRepo.DeleteIfExists(childNode.Id)

			if nil != childNode.Container {
				this.containerProvider.DeleteContainerIfExists(
					childNode.Id,
				)
			}
			defer waitGroup.Done()
		}(childNode)
	}

	waitGroup.Wait()

}
