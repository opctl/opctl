package local

import (
	"os"
	"path/filepath"
)

func (np nodeProvider) KillNodeIfExists(
	nodeID string,
) error {
	pIDOfLockOwner := np.lockfile.PIdOfOwner(
		filepath.Join(
			np.dataDir.Path(),
			"pid.lock",
		),
	)
	if pIDOfLockOwner != 0 {
		nodeProcess, err := os.FindProcess(pIDOfLockOwner)
		if err != nil {
			return err
		}

		if nodeProcess != nil {
			return nodeProcess.Kill()
		}
	}

	return nil
}
