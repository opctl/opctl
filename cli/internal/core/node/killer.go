package node

import (
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// Killer exposes the "node kill" sub command
type Killer interface {
	Kill() error
}

// newKiller returns an initialized "node kill" command
func newKiller(
	nodeProvider nodeprovider.NodeProvider,
) Killer {
	return _killer{
		nodeProvider,
	}
}

type _killer struct {
	nodeProvider nodeprovider.NodeProvider
}

func (ivkr _killer) Kill() error {
	return ivkr.nodeProvider.KillNodeIfExists("")
}
