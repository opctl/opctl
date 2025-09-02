package lima

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/internal/unsudo"
)

func limaBinPath(
	dataDir string,
) string {
	return filepath.Join(dataDir, "vms", "limactl")
}

// inflateLimaBinaries extracts embedded binaries for the current platform
func inflateLimaBinaries(binaryPath string) error {

	// Check if binary already exists and is executable
	if info, err := os.Stat(binaryPath); err == nil && info.Mode().Perm()&0111 != 0 {
		return nil
	}

	// Write limactl binary
	if err := writeBinary(binaryPath, limactlBytes); err != nil {
		return fmt.Errorf("failed to write limactl binary: %w", err)
	}

	// Write lima-guestagent binary
	guestagentPath := filepath.Join(filepath.Dir(binaryPath), "lima-guestagent.Linux-aarch64")
	if err := writeBinary(guestagentPath, limaLinuxGuestAgentBytes); err != nil {
		return fmt.Errorf("failed to write guest agent binary: %w", err)
	}

	return nil
}

// writeBinary writes the binary data to the specified path and makes it executable
func writeBinary(path string, data []byte) error {
	if err := unsudo.CreateFile(path, data); err != nil {
		return fmt.Errorf("failed to create binary file: %w", err)
	}
	return nil
}
