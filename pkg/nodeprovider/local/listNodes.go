package local

import (
	"fmt"
	"github.com/opspec-io/opctl/pkg/node"
)

func (this nodeProvider) ListNodes() (nodes []*node.InfoView, err error) {
	if nodeProcessId := this.nodeRepo.GetIfExists(); 0 != nodeProcessId {

		if this.psCanary.IsAlive(nodeProcessId) {
			nodes = []*node.InfoView{{}}
			return
		}

		fmt.Println("Found previous node process no longer running; cleaning up")
		this.nodeRepo.DeleteIfExists()
	}

	return
}
