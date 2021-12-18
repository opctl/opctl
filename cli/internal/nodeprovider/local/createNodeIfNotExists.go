//go:build darwin || dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/opctl/opctl/sdks/go/node"
)

func (np nodeProvider) CreateNodeIfNotExists(ctx context.Context) (node.Node, error) {
	apiClientNode, err := newAPIClientNode(np.config.ListenAddress)
	if err != nil {
		return nil, err
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
		np.config.ListenAddress,
		"--container-runtime",
		np.config.ContainerRuntime,
		"node",
		"create",
	)

	// don't inherit env; some things like jenkins track and kill processes via injecting env vars
	nodeCmd.Env = []string{
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	}

	// ensure node gets it's own process group
	nodeCmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	nodeLogFilePath := filepath.Join(np.dataDir.Path(), "node.log")
	nodeLogFile, err := os.OpenFile(nodeLogFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	nodeCmd.Stderr = io.MultiWriter(nodeLogFile, os.Stderr)
	nodeCmd.Stdout = io.MultiWriter(nodeLogFile, os.Stdout)

	if err := nodeCmd.Start(); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if err := apiClientNode.Liveness(ctx); err == nil {
					cancel()
				}
				time.Sleep(time.Second)
			}
		}
	}()

	<-ctx.Done()

	if err != nil {
		return nil, fmt.Errorf("failed to create daemonized opctl node: %w", err)
	}

	return apiClientNode, nil
}
