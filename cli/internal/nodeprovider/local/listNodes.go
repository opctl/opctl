package local

import (
	"path/filepath"

	"github.com/opctl/opctl/cli/internal/nodeprovider"
)

func (np nodeProvider) ListNodes() ([]nodeprovider.NodeHandle, error) {
	pIDOfLockOwner := np.lockfile.PIdOfOwner(
		filepath.Join(
			np.dataDir.Path(),
			"pid.lock",
		),
	)
	if 0 != pIDOfLockOwner {
		nodeHandle, err := newNodeHandle(np.listenAddress)
		if nil != err {
			return nil, err
		}

		return []nodeprovider.NodeHandle{
			nodeHandle,
		}, nil
	}

	return []nodeprovider.NodeHandle{}, nil
}
