package local

import (
	"context"
	"syscall"

	"github.com/opctl/opctl/cli/internal/pidfile"
)

func (np nodeProvider) KillNodeIfExists(
	ctx context.Context,
) error {
	nodeProcess, err := pidfile.TryGetProcess(
		ctx,
		np.config.DataDir,
	)
	if err != nil {
		return err
	}
	if nodeProcess == nil {
		return nil
	}

	return nodeProcess.SendSignalWithContext(ctx, syscall.SIGTERM)
}
