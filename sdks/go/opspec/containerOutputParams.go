package opspec

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/opctl/opctl/sdks/go/model"
)

var containerOutputName = regexp.MustCompile(`^[-_a-zA-Z0-9]+$`)

// ValidateContainerOutputParams checks declarations without interpreting their data.
func ValidateContainerOutputParams(declarations map[string]*model.ParamSpec, bindings map[string]string) error {
	if declarations == nil {
		return nil
	}
	if len(declarations) == 0 || len(bindings) == 0 {
		return fmt.Errorf("container outputParams requires declarations and nonempty outputs bindings")
	}
	names := make([]string, 0, len(declarations))
	for name := range declarations {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		p := declarations[name]
		if !containerOutputName.MatchString(name) || p == nil {
			return fmt.Errorf("invalid container outputParams declaration %q", name)
		}
		count := 0
		for _, present := range []bool{p.Array != nil, p.Boolean != nil, p.Dir != nil, p.File != nil, p.Number != nil, p.Object != nil, p.Socket != nil, p.String != nil} {
			if present {
				count++
			}
		}
		if count != 1 {
			return fmt.Errorf("container outputParams %q must declare exactly one type", name)
		}
		var fallback interface{}
		switch {
		case p.String != nil:
			fallback = p.String.Default
		case p.Number != nil:
			fallback = p.Number.Default
		case p.Boolean != nil:
			fallback = p.Boolean.Default
		case p.Array != nil:
			fallback = p.Array.Default
		case p.Object != nil:
			fallback = p.Object.Default
		default:
			return fmt.Errorf("container outputParams %q supports only string, number, boolean, array, and object; use mount/socket bindings for filesystem values", name)
		}
		if p.DefaultPresent || fallback != nil {
			return fmt.Errorf("container outputParams %q cannot have a default; every declared member must be produced", name)
		}
	}
	return nil
}
