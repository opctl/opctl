package opspec

import "github.com/opctl/opctl/sdks/go/model"

// HasContainerResults checks the authored graph, without resolving op references.
func HasContainerResults(call *model.CallSpec) bool {
	if call == nil {
		return false
	}
	if call.Container != nil && (call.Container.Stdout != "" || call.Container.Stderr != "" || len(call.Container.Outputs) > 0) {
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
