package container

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
)

var resultBinding = regexp.MustCompile(`^\$\([-_a-zA-Z0-9]+\)$`)

func validateResultBindings(spec *model.ContainerCallSpec) error {
	seen := map[string]bool{}
	targets := []string{spec.Stdout, spec.Stderr}
	if spec.Outputs != nil && len(spec.Outputs) == 0 {
		return fmt.Errorf("container outputs must have at least one named binding")
	}
	for name, ref := range spec.Outputs {
		if !regexp.MustCompile(`^[-_a-zA-Z0-9]+$`).MatchString(name) || !resultBinding.MatchString(ref) {
			return fmt.Errorf("invalid named container output binding")
		}
		targets = append(targets, ref)
	}
	for _, target := range targets {
		if target == "" {
			continue
		}
		if !resultBinding.MatchString(target) {
			return fmt.Errorf("container output target must be a whole identifier reference")
		}
		if seen[target] {
			return fmt.Errorf("duplicate container output target %s", target)
		}
		seen[target] = true
	}
	for _, mounts := range []map[string]interface{}{spec.Files, spec.Dirs} {
		for _, value := range mounts {
			if ref, ok := value.(string); ok && seen[ref] {
				return fmt.Errorf("duplicate container output target %s", ref)
			}
		}
	}
	for _, ref := range spec.Sockets {
		if seen[opspec.NameToRef(opspec.RefToName(ref))] {
			return fmt.Errorf("duplicate container output target %s", ref)
		}
	}
	return nil
}

func validateOutputsMount(call *model.ContainerCall) error {
	if _, exists := call.EnvVars["OPCTL_OUTPUTS"]; exists {
		return fmt.Errorf("OPCTL_OUTPUTS is reserved when container outputs are requested")
	}
	for _, mounts := range []map[string]string{call.Dirs, call.Files, call.Sockets} {
		for destination := range mounts {
			clean := path.Clean(destination)
			if clean == "/" || clean == "/opctl" || clean == "/opctl/outputs" || strings.HasPrefix(clean, "/opctl/outputs/") {
				return fmt.Errorf("container mount overlaps reserved /opctl/outputs directory")
			}
		}
	}
	return nil
}
