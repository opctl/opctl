//go:build darwin || dragonfly || freebsd || linux || nacl || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package local

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/opctl/opctl/sdks/go/node"
)

func (np nodeProvider) CreateNodeIfNotExists(
	ctx context.Context,
) (node.Node, error) {
	apiClientNode, err := newAPIClientNode(np.config.APIListenAddress)
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
		cmdArgs = append([]string{"-n", pathToOpctlBin}, cmdArgs...)
	}

	cmd := exec.Command(
		cmdName,
		cmdArgs...,
	)

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

	// always try to create a node; this avoids races
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// try to connect to existing opctl node
	err = apiClientNode.Liveness(ctx)
	if err == nil {
		return apiClientNode, nil
	}

	if os.Geteuid() != 0 {
		return nil, errors.New("re-run command with sudo")
	}

	return nil, fmt.Errorf("failed to create daemonized opctl node")
}
