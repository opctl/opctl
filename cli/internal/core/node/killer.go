package node

import (
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// Killer exposes the "node kill" sub command
type Killer interface {
	Kill()
}

// newKiller returns an initialized "node kill" command
func newKiller(
	cliExiter cliexiter.CliExiter,
	nodeProvider nodeprovider.NodeProvider,
) Killer {
	return _killer{
		cliExiter,
		nodeProvider,
	}
}

type _killer struct {
	cliExiter    cliexiter.CliExiter
	nodeProvider nodeprovider.NodeProvider
}

func (ivkr _killer) Kill() {
	err := ivkr.nodeProvider.KillNodeIfExists("")
	if nil != err {
		ivkr.cliExiter.Exit(cliexiter.ExitReq{Message: err.Error(), Code: 1})
		return // support fake exiter
	}
}
