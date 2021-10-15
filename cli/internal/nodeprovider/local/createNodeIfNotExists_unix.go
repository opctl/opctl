//go:build darwin || dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package local

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/opctl/opctl/sdks/go/node"
)

func (np nodeProvider) CreateNodeIfNotExists(ctx context.Context) (node.Node, error) {
	nodes, err := np.ListNodes()
	if err != nil {
		return nil, err
	}

	apiClientNode, err := newAPIClientNode(np.listenAddress)
	if err != nil {
		return nil, err
	}

	if len(nodes) > 0 {
		return apiClientNode, nil
	}

	pathToOpctlBin, err := os.Executable()
	if err != nil {
		return nil, err
	}

	pathToOpctlBin, err = filepath.EvalSymlinks(pathToOpctlBin)
	if err != nil {
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
	if err != nil {
		return nil, err
	}

	nodeCmd.Stderr = nodeLogFile
	nodeCmd.Stdout = nodeLogFile

	if err := nodeCmd.Start(); err != nil {
		return nil, err
	}

	err = apiClientNode.Liveness(ctx)
	nodeLogBytes, _ := os.ReadFile(nodeLogFilePath)
	fmt.Println(string(nodeLogBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create daemonized opctl node: %w", err)
	}

	return apiClientNode, nil
}
