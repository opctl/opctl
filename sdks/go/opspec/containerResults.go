package opspec

import (
	"fmt"
	"github.com/opctl/opctl/sdks/go/model"
)

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
	version, _ := SchemaRequirement(call)
	return version
}

// SchemaRequirement identifies the highest required schema and its first call path.
func SchemaRequirement(call *model.CallSpec) (int, string) {
	return schemaRequirement(call, "run")
}
func schemaRequirement(call *model.CallSpec, path string) (int, string) {
	version, field := 1, ""
	if call == nil {
		return version, field
	}
	if c := call.Container; c != nil {
		if c.Stdout != "" || c.Stderr != "" {
			version, field = 2, path+".container.stdout/stderr"
		}
		if len(c.Outputs) > 0 {
			version, field = 3, path+".container.outputs"
		}
		if c.OutputParams != nil {
			version, field = 4, path+".container.outputParams"
		}
	}
	for i, group := range []*[]*model.CallSpec{call.Serial, call.Parallel} {
		if group != nil {
			for j, child := range *group {
				needed, childField := schemaRequirement(child, fmt.Sprintf("%s.%s[%d]", path, []string{"serial", "parallel"}[i], j))
				if needed > version {
					version, field = needed, childField
				}
			}
		}
	}
	for i, loop := range []*model.CallSpec{serialLoopRun(call), parallelLoopRun(call)} {
		if loop != nil {
			needed, childField := schemaRequirement(loop, path+"."+[]string{"serialLoop", "parallelLoop"}[i]+".run")
			if needed > version {
				version, field = needed, childField
			}
		}
	}
	return version, field
}
func serialLoopRun(call *model.CallSpec) *model.CallSpec {
	if call.SerialLoop != nil {
		return &call.SerialLoop.Run
	}
	return nil
}
func parallelLoopRun(call *model.CallSpec) *model.CallSpec {
	if call.ParallelLoop != nil {
		return &call.ParallelLoop.Run
	}
	return nil
}
