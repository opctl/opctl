package local

import (
	"github.com/opctl/opctl/node"
)

func (np nodeProvider) ListNodes() (nodes []*node.InfoView, err error) {
	pIdOfLockOwner := np.lockfile.PIdOfOwner(np.lockFilePath)
	if 0 != pIdOfLockOwner {
		nodes = []*node.InfoView{{}}
	}
	return
}
