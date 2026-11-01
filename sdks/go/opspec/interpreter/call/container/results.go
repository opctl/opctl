package container

import (
	"fmt"
	"regexp"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
)

var resultBinding = regexp.MustCompile(`^\$\([-_a-zA-Z0-9]+\)$`)

func validateResultBindings(spec *model.ContainerCallSpec) error {
	seen := map[string]bool{}
	for _, target := range []string{spec.Stdout, spec.Stderr} {
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
