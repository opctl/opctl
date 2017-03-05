package core

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"sync"
)

func (this _core) KillOp(
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
