package core

//go:generate counterfeiter -o ./fakeCallKiller.go --fake-name fakeCallKiller ./ callKiller

import (
	"sync"

	"github.com/opctl/sdk-golang/model"

	"github.com/opctl/sdk-golang/node/core/containerruntime"
)

type callKiller interface {
	Kill(callID string)
}

func newCallKiller(
	callStore callStore,
	containerRuntime containerruntime.ContainerRuntime,
) callKiller {
	return _callKiller{
		callStore:        callStore,
		containerRuntime: containerRuntime,
	}
}

type _callKiller struct {
	callStore        callStore
	containerRuntime containerruntime.ContainerRuntime
}

func (this _callKiller) Kill(
	callID string,
) {
	this.callStore.SetIsKilled(callID)
	this.containerRuntime.DeleteContainerIfExists(callID)

	var waitGroup sync.WaitGroup

	for _, childCallGraph := range this.callStore.ListWithParentID(callID) {
		// recursively kill all child calls
		waitGroup.Add(1)
		go func(childCallGraph *model.DCG) {
			defer waitGroup.Done()
			this.Kill(childCallGraph.Id)
		}(childCallGraph)
	}

	waitGroup.Wait()

}
