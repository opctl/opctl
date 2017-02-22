package local

import "fmt"

func (this nodeProvider) KillNodeIfExists(
	nodeId string,
) (err error) {
	defer this.nodeRepo.DeleteIfExists()
	if nodeProcessId := this.nodeRepo.GetIfExists(); 0 != nodeProcessId {

		nodeProcess, findErr := this.os.FindProcess(nodeProcessId)
		if nil != findErr {
			fmt.Printf("error while killing node: findErr was: %v\n", findErr)
		}
		if nil != nodeProcess {
			killErr := nodeProcess.Kill()
			if nil != killErr {
				fmt.Printf("error while killing node: killErr was: %v\n", killErr)
			}

			fmt.Printf("killed node w/ PID: %v\n", nodeProcess.Pid)
		}
	}
	return
}
