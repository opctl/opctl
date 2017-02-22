package local

import (
	"fmt"
	"github.com/opspec-io/opctl/pkg/node"
	"strconv"
	"syscall"
)

func (this nodeProvider) ListNodes() (nodes []*node.InfoView, err error) {
	if nodeProcessId := this.nodeRepo.GetIfExists(); 0 != nodeProcessId {

		nodeProcess, findErr := this.os.FindProcess(nodeProcessId)
		if nil != findErr {
			fmt.Printf("error while listing nodes: findErr was: %v\n", findErr)
		}

		if nil != findErr || nil != nodeProcess.Signal(syscall.Signal(0)) {
			// cleanup
			this.KillNodeIfExists(strconv.Itoa(nodeProcessId))
		} else {
			// can be found && can be signalled w/ out error
			nodes = []*node.InfoView{{}}
		}

	}

	return
}
