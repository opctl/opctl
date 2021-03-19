package loop

import (
	"github.com/opctl/opctl/sdks/go/model"
)

// DeScope de-scopes loop vars (index, key, value)
func DeScope(
	parentScope map[string]*model.Value,
	callSpecLoopRange interface{},
	loopVarsSpec *model.LoopVarsSpec,
	iterationScope map[string]*model.Value,
) map[string]*model.Value {
	if loopVarsSpec == nil {
		return parentScope
	}

	outboundScope := map[string]*model.Value{}
	for varName, varData := range iterationScope {
		outboundScope[varName] = varData
	}

	// restore vars shadowed in the loop
	if loopVarsSpec.Index != nil {
		outboundScope[*loopVarsSpec.Index] = parentScope[*loopVarsSpec.Index]
	}

	if callSpecLoopRange != nil {
		return outboundScope
	}

	if loopVarsSpec.Key != nil {
		outboundScope[*loopVarsSpec.Key] = parentScope[*loopVarsSpec.Key]
	}
	if loopVarsSpec.Value != nil {
		outboundScope[*loopVarsSpec.Value] = parentScope[*loopVarsSpec.Value]
	}

	return outboundScope
}
