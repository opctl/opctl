// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package local

import (
	"github.com/opctl/opctl/cli/internal/model"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

func (np nodeProvider) CreateNode() (nodeInfo *model.NodeInfoView, err error) {
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
		Setpgid: true,
	}

	err = nodeCmd.Start()
	if nil != err {
		panic(err)
	}

	nodeInfo = &model.NodeInfoView{}

	return
}
