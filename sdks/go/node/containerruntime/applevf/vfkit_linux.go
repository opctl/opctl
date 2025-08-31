//go:build linux
// +build linux

package applevf

import (
	"fmt"

	"github.com/opctl/opctl/sdks/go/node/containerruntime"
	"github.com/opctl/opctl/sdks/go/node/containerruntime/fakes"
)

func New(
	cacheDir string,
) (containerruntime.ContainerRuntime, error) {
	// vfkit is not supported on linux
	return &fakes.FakeContainerRuntime{}, fmt.Errorf("appleVF container runtime is not supported on linux")
}
