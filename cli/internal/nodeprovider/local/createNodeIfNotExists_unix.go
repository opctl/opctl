// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package local

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/opctl/opctl/sdks/go/node"
)

func (np nodeProvider) CreateNodeIfNotExists() (node.Node, error) {
	nodes, err := np.ListNodes()
	if nil != err {
		return nil, err
	}

	apiClientNode, err := newAPIClientNode(np.listenAddress)
	if nil != err {
		return nil, err
	}

	if len(nodes) > 0 {
		return apiClientNode, nil
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
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
	}

	// ensure node gets it's own process group
	nodeCmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
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

	err = apiClientNode.Liveness(context.TODO())
	nodeLogBytes, _ := ioutil.ReadFile(nodeLogFilePath)
	fmt.Println(string(nodeLogBytes))
	if nil != err {
		return nil, fmt.Errorf("Error encountered creating daemonized opctl node; error was %s", err)
	}

	return apiClientNode, nil
}
