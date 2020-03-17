package local

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/opctl/opctl/cli/internal/datadir"
	"github.com/opctl/opctl/cli/internal/model"
)

func (np nodeProvider) CreateNodeIfNotExists() (model.NodeHandle, error) {
	nodes, err := np.ListNodes()
	if nil != err {
		return nil, err
	}

	nodeHandle, err := newNodeHandle()
	if nil != err {
		return nil, err
	}

	if len(nodes) > 0 {
		return nodeHandle, nil
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
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	if err := nodeCmd.Start(); nil != err {
		return nil, err
	}

	dataDir, err := datadir.New(nil)
	if nil != err {
		return nil, err
	}

	if err := nodeHandle.APIClient().Liveness(context.TODO()); nil != err {
		// Because the command is exec'd as a parentless process, error's aren't available to us.
		// To work around this the 'node create' call writes errors to the data dir, which we obtain here and error if found.
		return nil, fmt.Errorf("Error encountered creating opctl daemon; error was: %v", dataDir.TryGetNodeCreateError())
	}

	return nodeHandle, nil
}
