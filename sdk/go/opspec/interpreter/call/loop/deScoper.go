package loop

//go:generate counterfeiter -o ./fakeDeScoper.go --fake-name FakeDeScoper ./ DeScoper

import (
	"github.com/opctl/opctl/sdk/go/model"
)

type DeScoper interface {
	// DeScope de-scopes loop vars (index, key, value)
	DeScope(
		parentScope map[string]*model.Value,
		scgLoop *model.SCGLoop,
		iterationScope map[string]*model.Value,
	) map[string]*model.Value
}

func NewDeScoper() DeScoper {
	return _deScoper{}
}

type _deScoper struct{}

func (ds _deScoper) DeScope(
	parentScope map[string]*model.Value,
	scgLoop *model.SCGLoop,
	iterationScope map[string]*model.Value,
) map[string]*model.Value {
	outboundScope := map[string]*model.Value{}
	for varName, varData := range iterationScope {
		outboundScope[varName] = varData
	}

	// restore vars shadowed in the loop
	if nil != scgLoop.Index {
		outboundScope[*scgLoop.Index] = parentScope[*scgLoop.Index]
	}

	if nil != scgLoop.For {
		return outboundScope
	}

	if nil != scgLoop.For && nil != scgLoop.For.Key {
		outboundScope[*scgLoop.For.Key] = parentScope[*scgLoop.For.Key]
	}
	if nil != scgLoop.For && nil != scgLoop.For.Value {
		outboundScope[*scgLoop.For.Value] = parentScope[*scgLoop.For.Value]
	}

	return outboundScope
}
