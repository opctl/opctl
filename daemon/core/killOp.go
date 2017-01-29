package core

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"sync"
)

func (this _core) KillOp(
	req model.KillOpReq,
) {
	this.dcgNodeRepo.DeleteIfExists(req.OpGraphId)

	var waitGroup sync.WaitGroup

	for _, childNode := range this.dcgNodeRepo.ListWithOpGraphId(req.OpGraphId) {
		go func(childNode *dcgNodeDescriptor) {
			waitGroup.Add(1)
			this.dcgNodeRepo.DeleteIfExists(childNode.Id)

			if nil != childNode.Container {
				this.containerEngine.DeleteContainerIfExists(
					childNode.Id,
				)
			}
			defer waitGroup.Done()
		}(childNode)
	}

	waitGroup.Wait()

}
