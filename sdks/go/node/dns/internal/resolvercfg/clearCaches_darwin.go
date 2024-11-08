package resolvercfg

import (
	"context"
	"fmt"
	"os/exec"
)

func clearCaches(
	ctx context.Context,
) error {
	cmd := exec.CommandContext(
		ctx,
		"dscacheutil",
		"-flushcache",
	)

	if outputBytes, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to flush dns cache: %w, %s", err, string(outputBytes))
	}

	cmd = exec.CommandContext(
		ctx,
		"killall",
		"-HUP",
		"mDNSResponder",
	)

	if outputBytes, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to flush mdns cache: %w, %s", err, string(outputBytes))
	}

	return nil
}
