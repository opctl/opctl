//go:build darwin || dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package local

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

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

	cmdCtx, cancel := context.WithCancel(
		context.Background(),
	)

	cmdName := pathToOpctlBin
	cmdArgs := []string{
		"--data-dir",
		np.config.DataDir,
		"--listen-address",
		np.config.ListenAddress,
		"--container-runtime",
		np.config.ContainerRuntime,
		"node",
		"create",
	}

	if os.Geteuid() != 0 {
		cmdName = "sudo"
		cmdArgs = append([]string{pathToOpctlBin}, cmdArgs...)
	}

	cmd := exec.CommandContext(
		cmdCtx,
		cmdName,
		cmdArgs...,
	)

	// don't inherit env; some things like jenkins track and kill processes via injecting env vars
	cmd.Env = []string{
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	}

	// buffer output
	nodeCmdOutput := bytes.Buffer{}
	cmd.Stderr = &nodeCmdOutput
	cmd.Stdout = &nodeCmdOutput
	cmd.Stdin = nil

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// own process group
		Setpgid: true,
	}

	// always try to create a node; this avoids races
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	defer func() {
		// ensure resources always cleaned up
		if cmd.ProcessState != nil {
			cmd.Wait()
		}
	}()

	// try to connect to existing opctl node
	err = apiClientNode.Liveness(ctx)
	if err == nil {
		return apiClientNode, nil
	}

	cancel()

	cmd.Wait()

	if os.Geteuid() != 0 {
		return nil, errors.New("re-run command with sudo")
	}

	return nil, fmt.Errorf("failed to create daemonized opctl node: %s", nodeCmdOutput.String())
}
