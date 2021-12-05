package k8s

import (
	"fmt"
)

func constructPodName(
	containerID string,
) string {
	return fmt.Sprintf("opctl-%s", containerID)
}
