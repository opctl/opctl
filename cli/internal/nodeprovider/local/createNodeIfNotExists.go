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
	"strings"
	"syscall"

	"github.com/opctl/opctl/cli/internal/pidfile"
	"github.com/opctl/opctl/sdks/go/node"
)

func (np nodeProvider) CreateNodeIfNotExists(
	ctx context.Context,
) (node.Node, error) {
	apiClientNode, err := newAPIClientNode(np.config.APIListenAddress)
	if err != nil {
		return nil, err
	}

	nodeProcess, err := pidfile.TryGetProcess(
		ctx,
		np.config.DataDir,
	)
	if err != nil {
		return nil, err
	}

	if nodeProcess != nil {
		return apiClientNode, nil
	}

	// no process running, need to daemonize node...

	pathToOpctlBin, err := os.Executable()
	if err != nil {
		return nil, err
	}

	pathToOpctlBin, err = filepath.EvalSymlinks(pathToOpctlBin)
	if err != nil {
		return nil, err
	}

	cmdName := pathToOpctlBin
	cmdArgs := []string{
		"--api-listen-address",
		np.config.APIListenAddress,
		"--container-runtime",
		np.config.ContainerRuntime,
		"--data-dir",
		np.config.DataDir,
		"--dns-listen-address",
		np.config.DNSListenAddress,
		"node",
		"create",
	}

	if os.Geteuid() != 0 {
		cmdName = "sudo"
		cmdArgs = append(
			[]string{
				// non interactive
				"-n",
				pathToOpctlBin,
			},
			cmdArgs...,
		)
	}

	cmd := exec.Command(
		cmdName,
		cmdArgs...,
	)

	cmd.Stdout = os.Stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// don't inherit env; some things like jenkins track and kill processes via injecting env vars
	cmd.Env = []string{
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		// set by sudo; passthru so we maintain provenance for use by "unsudo"
		fmt.Sprintf("SUDO_GID=%s", os.Getenv("SUDO_GID")),
		fmt.Sprintf("SUDO_UID=%s", os.Getenv("SUDO_UID")),
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		// own process group
		Setpgid: true,
	}

	ctx, cancel := context.WithCancel(ctx)

	var daemonErr error
	go func() {
		daemonErr = cmd.Run()
		defer cancel()
	}()

	var timeoutErr error
	go func() {
		timeoutErr = apiClientNode.Liveness(ctx)
		defer cancel()
	}()

	<-ctx.Done()

	// stop buffering Stderr
	cmd.Stderr = os.Stderr

	// handle error daemonizing
	if _, ok := daemonErr.(*exec.ExitError); ok {
		// handle race
		if strings.Contains(stderr.String(), "node already running") {
			return apiClientNode, nil
		}

		// handle "sudo: a password is required"
		if strings.Contains(stderr.String(), "password") {
			return nil, fmt.Errorf("re-run command with sudo")
		}

		return nil, errors.New(stderr.String())
	}

	// handle timeout reaching opctl API
	if timeoutErr != nil {
		return nil, fmt.Errorf(
			"timeout reaching opctl API at %s; try re-running command",
			np.config.APIListenAddress,
		)
	}

	return apiClientNode, nil
}
