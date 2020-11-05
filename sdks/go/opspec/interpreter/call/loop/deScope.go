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
	if nil == loopVarsSpec {
		return parentScope
	}

	outboundScope := map[string]*model.Value{}
	for varName, varData := range iterationScope {
		outboundScope[varName] = varData
	}

	// restore vars shadowed in the loop
	if nil != loopVarsSpec.Index {
		outboundScope[*loopVarsSpec.Index] = parentScope[*loopVarsSpec.Index]
	}

	if nil != callSpecLoopRange {
		return outboundScope
	}

	if nil != loopVarsSpec.Key {
		outboundScope[*loopVarsSpec.Key] = parentScope[*loopVarsSpec.Key]
	}
	if nil != loopVarsSpec.Value {
		outboundScope[*loopVarsSpec.Value] = parentScope[*loopVarsSpec.Value]
	}

	return outboundScope
}
