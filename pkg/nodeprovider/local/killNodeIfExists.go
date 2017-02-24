package local

import (
	"os"
)

func (this nodeProvider) KillNodeIfExists(
	nodeId string,
) (err error) {

	pIdOfLockOwner := this.lockfile.PIdOfOwner(lockFilePath())
	if 0 != pIdOfLockOwner {
		var nodeProcess *os.Process
		nodeProcess, err = this.os.FindProcess(pIdOfLockOwner)
		if nil != nodeProcess {
			err = nodeProcess.Kill()
		}
	}
	return
}
