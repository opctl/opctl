package docker

import "fmt"

// rootCallIDLabel holds the ID of the root op call a container runs in.
// do not change this label as it might break external consumers
const rootCallIDLabel = "opctl.rootCallId"

func getContainerName(opctlContainerID string) string {
	return fmt.Sprintf("%s%s", containerNamePrefix, opctlContainerID)
}
