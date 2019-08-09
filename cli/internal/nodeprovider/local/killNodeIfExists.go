package local

import (
	"os"
)

func (np nodeProvider) KillNodeIfExists(
	nodeId string,
) (err error) {

	pIdOfLockOwner := np.lockfile.PIdOfOwner(np.lockFilePath)
	if 0 != pIdOfLockOwner {
		var nodeProcess *os.Process
		nodeProcess, err = np.os.FindProcess(pIdOfLockOwner)
		if nil != nodeProcess {
			err = nodeProcess.Kill()
		}
	}
	return
}
