package core

import "github.com/opctl/opctl/node"

func (this _core) NodeCreate(
	opts NodeCreateOpts,
) {
	config := node.Config{}
	if "" != opts.DataDir {
		config.DataDir = &opts.DataDir
	}
	node.New(config)
}
