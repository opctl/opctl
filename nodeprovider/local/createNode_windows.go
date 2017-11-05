package local

import (
	"github.com/opctl/opctl/node"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func (np nodeProvider) CreateNode() (nodeInfo *node.InfoView, err error) {
	pathToOpctlBin, err := os.Executable()
	if nil != err {
		return nil, err
	}

	pathToOpctlBin, err = filepath.EvalSymlinks(pathToOpctlBin)
	if nil != err {
		return nil, err
	}

	nodeCmd := exec.Command(
		pathToOpctlBin,
		"node",
		"create",
	)

	nodeCmd.SysProcAttr = &syscall.SysProcAttr{
		// ensure node gets it's own process group
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	err = nodeCmd.Start()
	if nil != err {
		panic(err)
	}

	nodeInfo = &node.InfoView{}

	return
}
