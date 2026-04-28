package resolvercfg

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Delete modifications to the current system
func Delete(
	ctx context.Context,
) error {
	if err := deleteScutilEntries(ctx); err != nil {
		return err
	}

	// Clean up legacy /etc/resolver/opctl_* files from older versions
	deleteLegacyResolverFiles()

	return clearCaches(ctx)
}

func deleteScutilEntries(ctx context.Context) error {
	// List all opctl resolver keys in the Dynamic Store
	listCmd := exec.CommandContext(ctx, "scutil")
	listCmd.Stdin = strings.NewReader(
		fmt.Sprintf("list %s[^/]+/DNS\n", scutilKeyPrefix),
	)

	output, err := listCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to list resolver configs via scutil: %w, %s", err, string(output))
	}

	// Parse output for keys to remove
	var removeCmds strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// scutil list output lines look like: "  subKey [0] = State:/Network/Service/opctl_foo/DNS"
		if idx := strings.Index(line, "= "); idx >= 0 {
			key := strings.TrimSpace(line[idx+2:])
			if strings.HasPrefix(key, scutilKeyPrefix) {
				removeCmds.WriteString(fmt.Sprintf("remove %s\n", key))
			}
		}
	}

	if removeCmds.Len() == 0 {
		return nil
	}

	removeCmd := exec.CommandContext(ctx, "scutil")
	removeCmd.Stdin = strings.NewReader(removeCmds.String())

	outputBytes, err := removeCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove resolver configs via scutil: %w, %s", err, string(outputBytes))
	}

	// scutil exits 0 even on permission errors, so check output
	if strings.Contains(string(outputBytes), "Permission denied") {
		return fmt.Errorf("failed to remove resolver configs via scutil: permission denied (run as root)")
	}

	return nil
}

// deleteLegacyResolverFiles removes /etc/resolver/opctl_* files left by older versions.
func deleteLegacyResolverFiles() {
	legacyDir := "/etc/resolver"
	legacyPrefix := "opctl_"

	_ = filepath.WalkDir(
		legacyDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil || d == nil {
				return nil
			}
			if strings.HasPrefix(d.Name(), legacyPrefix) {
				_ = os.Remove(path)
			}
			return nil
		},
	)
}
