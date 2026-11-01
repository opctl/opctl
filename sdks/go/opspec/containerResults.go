package opspec

import "github.com/opctl/opctl/sdks/go/model"

// HasContainerResults checks the authored graph, without resolving op references.
func HasContainerResults(call *model.CallSpec) bool {
	if call == nil {
		return false
	}
	if call.Container != nil && (call.Container.Stdout != "" || call.Container.Stderr != "" || len(call.Container.Outputs) > 0 || call.Container.OutputParams != nil) {
		return true
	}
	for _, group := range []*[]*model.CallSpec{call.Serial, call.Parallel} {
		if group != nil {
			for _, child := range *group {
				if HasContainerResults(child) {
					return true
				}
			}
		}
	}
	if call.SerialLoop != nil && HasContainerResults(&call.SerialLoop.Run) {
		return true
	}
	return call.ParallelLoop != nil && HasContainerResults(&call.ParallelLoop.Run)
}

// MinimumSchema reports graph features independently of a bootstrap envelope label.
func MinimumSchema(call *model.CallSpec) int {
	version := 1
	if call == nil {
		return version
	}
	if c := call.Container; c != nil {
		if c.Stdout != "" || c.Stderr != "" {
			version = 2
		}
		if len(c.Outputs) > 0 {
			version = 3
		}
		if c.OutputParams != nil {
			version = 4
		}
	}
	for _, group := range []*[]*model.CallSpec{call.Serial, call.Parallel} {
		if group != nil {
			for _, child := range *group {
				if needed := MinimumSchema(child); needed > version {
					version = needed
				}
			}
		}
	}
	if call.SerialLoop != nil {
		if needed := MinimumSchema(&call.SerialLoop.Run); needed > version {
			version = needed
		}
	}
	if call.ParallelLoop != nil {
		if needed := MinimumSchema(&call.ParallelLoop.Run); needed > version {
			version = needed
		}
	}
	return version
}
