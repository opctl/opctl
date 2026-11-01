package goop

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
)

// ContainerSpec prepares a bootstrap and returns the host dir to bind as
// the module scope value. Local replace targets retain their host paths.
func sourceSpec(opDir, phase string) (*model.ContainerCallSpec, string, error) {
	return sourceSpecInRepository(opDir, phase, "")
}

// sourceSpecInRepository bounds git-fetched Go sources and local replacements
// to their resolved repository root. An empty root retains local-op behavior.
func sourceSpecInRepository(opDir, phase, repositoryRoot string) (*model.ContainerCallSpec, string, error) {
	if phase != "plan" && phase != "describe" {
		return nil, "", fmt.Errorf("unsupported bootstrap phase %q", phase)
	}
	source, err := Discover(opDir, repositoryRoot)
	if err != nil {
		return nil, "", err
	}
	return source.ContainerSpec(phase), source.ModuleRoot, nil
}
