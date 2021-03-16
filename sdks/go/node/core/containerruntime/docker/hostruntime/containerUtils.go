package hostruntime

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

type containerUtils struct {
	inAContainer func() bool
	getDockerID  func() (string, error)
}

var defaultContainerUtils = containerUtils{
	inAContainer,
	getDockerID,
}

func inAContainer() bool {
	// Only linux containers are supported.
	// If GOOS is not linux we assume we're not in a container.
	if runtime.GOOS != "linux" {
		return false
	}

	// https://github.com/testcontainers/testcontainers-go/blob/f4135ba4793efb7609477e94715a74e111ddf67b/docker.go#L600
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	return false
}

func getDockerID() (string, error) {
	// Only linux containers are supported.
	if runtime.GOOS != "linux" {
		return "", fmt.Errorf("GOOS %s not supported", runtime.GOOS)
	}

	// https://forums.docker.com/t/get-a-containers-full-id-from-inside-of-itself/37237
	cpuset, err := os.ReadFile("/proc/1/cpuset")
	if err != nil {
		return "", err
	}

	tokens := strings.Split(string(cpuset), "/")
	dockerID := tokens[len(tokens)-1][0:13]

	return dockerID, nil
}
