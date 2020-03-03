package node

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

import (
	"github.com/opctl/opctl/cli/internal/cliexiter"
	"github.com/opctl/opctl/cli/internal/core/node/creater"
	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

// Node exposes the "node" sub command
//counterfeiter:generate -o fakes/node.go . Node
type Node interface {
	creater.Creater
	Killer
}

// New returns an initialized "node" sub command
func New(
	cliExiter cliexiter.CliExiter,
	nodeProvider nodeprovider.NodeProvider,
) Node {
	return _node{
		Creater: creater.New(),
		Killer: newKiller(
			cliExiter,
			nodeProvider,
		),
	}
}

type _node struct {
	creater.Creater
	Killer
}
