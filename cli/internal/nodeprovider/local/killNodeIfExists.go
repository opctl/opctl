package local

import (
	"os"
	"path/filepath"
	"strings"
)

func (np nodeProvider) KillNodeIfExists() error {
	pID, err := getPIDFromFile(
		filepath.Join(
			np.config.DataDir,
			"pid.lock",
		),
	)
	if err != nil {
		if os.IsNotExist(err) {
			// already killed or our mutex was manually removed
			return nil
		}

		return err
	}

	nodeProcess, err := os.FindProcess(pID)
	if err != nil {
		return err
	}

	if nodeProcess != nil {
		err = nodeProcess.Kill()
		if nil != err && !strings.Contains(err.Error(), "os: process already finished") {
			return err
		}

		// ignore errors because we don't care; we just need it to have exited
		nodeProcess.Wait()
	}

	return nil
}
