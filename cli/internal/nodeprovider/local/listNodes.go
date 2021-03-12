package local

import (
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/node"
)

func (np nodeProvider) ListNodes() ([]node.Node, error) {
	pIDOfLockOwner := np.lockfile.PIdOfOwner(
		filepath.Join(
			np.dataDir.Path(),
			"pid.lock",
		),
	)
	if 0 != pIDOfLockOwner {
		apiClientNode, err := newAPIClientNode(np.listenAddress)
		if nil != err {
			return nil, err
		}

		return []node.Node{
			apiClientNode,
		}, nil
	}

	return []node.Node{}, nil
}
