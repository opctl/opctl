package docker

import "fmt"

func getContainerName(opctlContainerID string) string {
	return fmt.Sprintf("opctl_%s", opctlContainerID)
}
