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
	if 0 != pIDOfLockOwner {
		nodeProcess, err := os.FindProcess(pIDOfLockOwner)
		if nil != err {
			return err
		}

		if nil != nodeProcess {
			return nodeProcess.Kill()
		}
	}

	return nil
}
