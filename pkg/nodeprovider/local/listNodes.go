package local

import (
	"github.com/opspec-io/opctl/pkg/node"
)

func (this nodeProvider) ListNodes() (nodes []*node.InfoView, err error) {
	pIdOfLockOwner := this.lockfile.PIdOfOwner(lockFilePath())
	if 0 != pIdOfLockOwner {
		nodes = []*node.InfoView{{}}
	}
	return
}
