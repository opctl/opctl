// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package local

import (
	"fmt"
	"github.com/opspec-io/opctl/pkg/node"
	"os/exec"
	"syscall"
)

func (this nodeProvider) CreateNode() (nodeInfo *node.InfoView, err error) {
	nodeCmd := exec.Command(
		"opctl",
		"node",
		"create",
	)

	nodeCmd.SysProcAttr = &syscall.SysProcAttr{
		// ensure node gets it's own process group
		Setpgid: true,
	}

	err = nodeCmd.Start()
	if nil != err {
		panic(err)
	}

	fmt.Printf("created node w/ PID: %v\n", nodeCmd.Process.Pid)

	this.nodeRepo.Add(nodeCmd.Process.Pid)

	nodeInfo = &node.InfoView{}

	return
}
