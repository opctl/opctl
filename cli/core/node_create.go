package core

import "github.com/opctl/opctl/node"

func (this _core) NodeCreate() {
	this.cliOutput.Info("starting node; ReST API will be bound to 0.0.0.0:42224")
	node.New()
}
