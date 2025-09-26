package unsudo

import (
	"context"
	"os/exec"
	"syscall"
)

// NewCmd that will run as the user & group who ran sudo
func NewCmd(ctx context.Context, name string, arg ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, arg...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(getSudoUID()),
			Gid: uint32(getSudoGID()),
		},
	}

	return cmd
}
