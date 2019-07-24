package local

import (
	"github.com/opctl/opctl/cli/types"
)

func (np nodeProvider) ListNodes() (nodes []*types.NodeInfoView, err error) {
	pIdOfLockOwner := np.lockfile.PIdOfOwner(np.lockFilePath)
	if 0 != pIdOfLockOwner {
		nodes = []*types.NodeInfoView{{}}
	}
	return
}
