package apireachabilityensurer

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeAPIReachabilityEnsurer.go --fake-name Fake ./ APIReachabilityEnsurer

import (
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
	"github.com/opctl/opctl/cli/internal/nodeprovider/local"
	"time"
)

type APIReachabilityEnsurer interface {
	// ensures a node is reachable by starting one if none found running
	// panics if any error(s) encountered
	Ensure()
}

func New(
	cliExiter cliexiter.CliExiter,
) APIReachabilityEnsurer {
	return _apiReachabilityEnsurer{
		cliExiter:    cliExiter,
		nodeProvider: local.New(),
	}
}

type _apiReachabilityEnsurer struct {
	cliExiter    cliexiter.CliExiter
	nodeProvider nodeprovider.NodeProvider
}

func (nae _apiReachabilityEnsurer) Ensure() {
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
