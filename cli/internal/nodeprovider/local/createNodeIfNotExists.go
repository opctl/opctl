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

	"github.com/opctl/opctl/cli/internal/euid0"
	"github.com/opctl/opctl/sdks/go/node"
	"github.com/opctl/opctl/sdks/go/node/containerruntime"
)

func (np nodeProvider) CreateNodeIfNotExists(
	ctx context.Context,
) (node.Node, error) {
	apiClientNode, err := newAPIClientNode(np.config.APIListenAddress)
	if err != nil {
		return nil, err
	}

	// check if node API reachable
	if err := apiClientNode.Liveness(ctx); err == nil {
		return apiClientNode, nil
	}

	// node API unreachable, need to daemonize node...
	if err := euid0.Ensure(); err != nil {
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

	cmd := exec.Command(
		pathToOpctlBin,
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
	)

	cmd.Stdout = os.Stdout

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	cmd.Env = daemonEnv()

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

		return nil, errors.New(stderr.String())
	}

	// handle timeout reaching node API
	if timeoutErr != nil {
		return nil, fmt.Errorf(
			"timeout reaching node API at %s; try re-running command",
			np.config.APIListenAddress,
		)
	}

	return apiClientNode, nil
}

// daemonEnv builds the environment for the daemonized node process.
//
// We deliberately don't inherit the full env; some things (e.g. jenkins) track
// and kill processes via injected env vars. We pass through only what the node
// needs: HOME & PATH, the SUDO_* provenance used by "unsudo", and the proxy
// vars so the node (and the op containers it runs, see
// containerruntime.ProxyEnvVars) can reach the network on hosts whose only
// egress route is an HTTP/HTTPS forward proxy (e.g. CI runners).
func daemonEnv() []string {
	env := []string{
		fmt.Sprintf("HOME=%s", os.Getenv("HOME")),
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		// set by sudo; passthru so we maintain provenance for use by "unsudo"
		fmt.Sprintf("SUDO_GID=%s", os.Getenv("SUDO_GID")),
		fmt.Sprintf("SUDO_UID=%s", os.Getenv("SUDO_UID")),
	}

	// passing nil yields every proxy var present in this process's environment
	for name, value := range containerruntime.ProxyEnvVars(nil) {
		env = append(env, fmt.Sprintf("%s=%s", name, value))
	}

	return env
}
