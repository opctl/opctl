package local

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

func (np nodeProvider) CreateNodeIfNotExists() (nodeprovider.NodeHandle, error) {
	nodes, err := np.ListNodes()
	if nil != err {
		return nil, err
	}

	nodeHandle, err := newNodeHandle(np.listenAddress)
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
		"--data-dir",
		np.dataDir.Path(),
		"--listen-address",
		np.listenAddress,
		"node",
		"create",
	)

	// don't inherit env; some things like jenkins track and kill processes via injecting env vars
	nodeCmd.Env = []string{
		fmt.Sprintf("LOCALAPPDATA=%s", os.Getenv("LOCALAPPDATA")),
	}

	// ensure node gets it's own process group
	nodeCmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}

	nodeLogFilePath := filepath.Join(np.dataDir.Path(), "node.log")
	nodeLogFile, err := os.Create(nodeLogFilePath)
	if nil != err {
		return nil, err
	}

	nodeCmd.Stderr = nodeLogFile
	nodeCmd.Stdout = nodeLogFile

	if err := nodeCmd.Start(); nil != err {
		return nil, err
	}

	err = nodeHandle.APIClient().Liveness(context.TODO())
	nodeLogBytes, _ := ioutil.ReadFile(nodeLogFilePath)
	fmt.Println(string(nodeLogBytes))
	if nil != err {
		return nil, fmt.Errorf("Error encountered creating daemonized opctl node")
	}

	return nodeHandle, nil
}
