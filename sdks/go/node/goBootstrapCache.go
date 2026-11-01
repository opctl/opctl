package node

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
)

// goBootstrapCache is shared by API describes and all op callers in this node.
type goBootstrapCache struct {
	runtime    containerruntime.ContainerRuntime
	stateStore stateStore
}

// Registry diagnostics may contain credentials. Bootstrap reports controlled errors.
type discardBootstrapEvents struct{}

func (discardBootstrapEvents) Publish(model.Event) {}
