package local

import (
	"github.com/opctl/opctl/cli/types"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func (np nodeProvider) CreateNode() (nodeInfo *types.NodeInfoView, err error) {
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

	nodeInfo = &types.NodeInfoView{}

	return
}
