package local

import (
	"fmt"
	"github.com/opspec-io/opctl/pkg/node"
	"os/exec"
)

func (this nodeProvider) CreateNode() (nodeInfo *node.InfoView, err error) {
	nodeCmd := exec.Command(
		"opctl",
		"node",
		"create",
	)

	err = nodeCmd.Start()
	if nil != err {
		panic(err)
	}

	fmt.Printf("created node w/ PID: %v\n", nodeCmd.Process.Pid)

	this.nodeRepo.Add(nodeCmd.Process.Pid)

	nodeInfo = &node.InfoView{}

	return
}
