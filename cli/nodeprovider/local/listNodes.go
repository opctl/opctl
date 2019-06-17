package local

import (
	"github.com/opctl/opctl/cli/model"
)

func (np nodeProvider) ListNodes() (nodes []*model.NodeInfoView, err error) {
	pIdOfLockOwner := np.lockfile.PIdOfOwner(np.lockFilePath)
	if 0 != pIdOfLockOwner {
		nodes = []*model.NodeInfoView{{}}
	}
	return
}
