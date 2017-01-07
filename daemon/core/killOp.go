package core

import (
	"github.com/opspec-io/sdk-golang/pkg/model"
	"sync"
)

func (this _core) KillOp(
	req model.KillOpReq,
) {
	this.nodeRepo.deleteIfExists(req.OpGraphId)

	var waitGroup sync.WaitGroup

	for _, childNode := range this.nodeRepo.listWithOpGraphId(req.OpGraphId) {
		go func(childNode *nodeDescriptor) {
			waitGroup.Add(1)
			this.nodeRepo.deleteIfExists(childNode.Id)

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
