package docker

import "fmt"

func getContainerName(opctlContainerID string) string {
	return fmt.Sprintf("%s%s", containerNamePrefix, opctlContainerID)
}
