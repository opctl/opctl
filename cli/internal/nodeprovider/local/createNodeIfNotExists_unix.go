// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package local

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/model"
)

func (np nodeProvider) CreateNodeIfNotExists() (
	nodeHandle model.NodeHandle,
	err error,
) {
	nodes, err := np.ListNodes()
	if nil != err {
		return nil, err
	}

	if len(nodes) > 0 {
		return newNodeHandle()
	}

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
		return nil, err
	}

	dataDir, err := datadir.New(nil)
	if nil != err {
		return nil, err
	}

	// Because the command is exec'd as a parentless process, error's aren't available to us.
	// To work around this the 'node create' call writes errors to the data dir, which we poll here and panic if found.
	time.Sleep(3 * time.Second)
	nodeCreateErr := dataDir.TryGetNodeCreateError()
	if nil != nodeCreateErr {
		return nil, fmt.Errorf("Error encountered creating opctl daemon; error was: %v", nodeCreateErr)
	}

	return newNodeHandle()
}
