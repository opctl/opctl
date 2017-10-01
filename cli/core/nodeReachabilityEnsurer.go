package core

//go:generate counterfeiter -o ./fakeNodeReachabilityEnsurer.go --fake-name fakeNodeReachabilityEnsurer ./ nodeReachabilityEnsurer

import (
	"github.com/opctl/opctl/nodeprovider"
	"github.com/opctl/opctl/nodeprovider/local"
	"github.com/opctl/opctl/util/cliexiter"
	"time"
)

type nodeReachabilityEnsurer interface {
	// ensures a node is reachable by starting one if none found running
	// panics if any error(s) encountered
	EnsureNodeReachable()
}

func newNodeReachabilityEnsurer(
	cliExiter cliexiter.CliExiter,
) nodeReachabilityEnsurer {
	return _nodeReachabilityEnsurer{
		cliExiter:    cliExiter,
		nodeProvider: local.New(),
	}
}

type _nodeReachabilityEnsurer struct {
	cliExiter    cliexiter.CliExiter
	nodeProvider nodeprovider.NodeProvider
}

func (nae _nodeReachabilityEnsurer) EnsureNodeReachable() {
	nodes, err := nae.nodeProvider.ListNodes()
	if nil != err {
		nae.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}

	if len(nodes) < 1 {
		nae.nodeProvider.CreateNode()
		// sleep to let the opctl node start
		// @TODO: add exp backoff to SDK websocket client so we don't need this
		<-time.After(time.Second * 3)
	}
}
