package local

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

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

	// don't inherit env; some things like jenkins track and kill processes via injecting env vars
	nodeCmd.Env = []string{
		fmt.Sprintf("OPCTL_DATA_DIR=%v", np.dataDir.Path()),
	}

	// ensure node gets it's own process group
	nodeCmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	nodeCmdOutput := &bytes.Buffer{}
	nodeCmd.Stdout = nodeCmdOutput
	nodeCmd.Stderr = nodeCmdOutput

	if err := nodeCmd.Start(); nil != err {
		return nil, err
	}

	if err := nodeHandle.APIClient().Liveness(context.TODO()); nil != err {
		fmt.Print(string(nodeCmdOutput.Bytes()))
		return nil, fmt.Errorf("Error encountered creating opctl daemon")
	}

	fmt.Print(string(nodeCmdOutput.Bytes()))

	return nodeHandle, nil
}
