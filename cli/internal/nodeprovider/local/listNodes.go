package local

import (
	"github.com/opctl/opctl/cli/internal/model"
)

func (np nodeProvider) ListNodes() ([]model.NodeHandle, error) {
	pIDOfLockOwner := np.lockfile.PIdOfOwner(np.lockFilePath)
	if 0 != pIDOfLockOwner {
		nodeHandle, err := newNodeHandle()
		if nil != err {
			return nil, err
		}

		return []model.NodeHandle{
			nodeHandle,
		}, nil
	}

	return []model.NodeHandle{}, nil
}
