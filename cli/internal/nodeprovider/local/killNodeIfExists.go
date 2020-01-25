package local

func (np nodeProvider) KillNodeIfExists(
	nodeID string,
) error {
	pIDOfLockOwner := np.lockfile.PIdOfOwner(np.lockFilePath)
	if 0 != pIDOfLockOwner {
		nodeProcess, err := np.os.FindProcess(pIDOfLockOwner)
		if nil != err {
			return err
		}

		if nil != nodeProcess {
			return nodeProcess.Kill()
		}
	}

	return nil
}
