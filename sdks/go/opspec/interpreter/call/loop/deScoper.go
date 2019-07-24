package loop

//go:generate counterfeiter -o ./fakeDeScoper.go --fake-name FakeDeScoper ./ DeScoper

import (
	"github.com/opctl/opctl/sdks/go/types"
)

type DeScoper interface {
	// DeScope de-scopes loop vars (index, key, value)
	DeScope(
		parentScope map[string]*types.Value,
		scgLoopRange interface{},
		scgLoopVars *types.SCGLoopVars,
		iterationScope map[string]*types.Value,
	) map[string]*types.Value
}

func NewDeScoper() DeScoper {
	return _deScoper{}
}

type _deScoper struct{}

func (ds _deScoper) DeScope(
	parentScope map[string]*types.Value,
	scgLoopRange interface{},
	scgLoopVars *types.SCGLoopVars,
	iterationScope map[string]*types.Value,
) map[string]*types.Value {
	if nil == scgLoopVars {
		return parentScope
	}

	outboundScope := map[string]*types.Value{}
	for varName, varData := range iterationScope {
		outboundScope[varName] = varData
	}

	// restore vars shadowed in the loop
	if nil != scgLoopVars.Index {
		outboundScope[*scgLoopVars.Index] = parentScope[*scgLoopVars.Index]
	}

	if nil != scgLoopRange {
		return outboundScope
	}

	if nil != scgLoopVars.Key {
		outboundScope[*scgLoopVars.Key] = parentScope[*scgLoopVars.Key]
	}
	if nil != scgLoopVars.Value {
		outboundScope[*scgLoopVars.Value] = parentScope[*scgLoopVars.Value]
	}

	return outboundScope
}
