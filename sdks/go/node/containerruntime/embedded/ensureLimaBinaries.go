package embedded

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

	guestArch := ""
	if runtime.GOARCH == "amd64" {
		guestArch = "x86_64"
	} else if runtime.GOARCH == "arm64" {
		guestArch = "aarch64"
	} else {
		return fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}

	// Write lima-guestagent binary
	guestagentPath := filepath.Join(
		filepath.Dir(binaryPath),
		fmt.Sprintf("lima-guestagent.Linux-%s", guestArch),
	)
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
