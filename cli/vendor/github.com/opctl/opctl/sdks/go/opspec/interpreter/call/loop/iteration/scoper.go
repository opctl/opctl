package iteration

//go:generate counterfeiter -o ./fakeScoper.go --fake-name FakeScoper ./ Scoper

import (
	"sort"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"

	"github.com/opctl/opctl/sdks/go/types"
)

type Scoper interface {
	// Scope scopes loop iteration vars (index, key, value)
	Scope(
		index int,
		scope map[string]*types.Value,
		scgLoopRange interface{},
		scgLoopVars *types.SCGLoopVars,
		opHandle types.DataHandle,
	) (
		map[string]*types.Value,
		error,
	)
}

func NewScoper() Scoper {
	return _scoper{
		loopableInterpreter: loopable.NewInterpreter(),
		valueInterpreter:    value.NewInterpreter(),
	}
}

type _scoper struct {
	loopableInterpreter loopable.Interpreter
	valueInterpreter    value.Interpreter
}

func (lpr _scoper) sortMap(
	m map[string]interface{},
) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}

	sort.Strings(names) //sort keys alphabetically
	return names
}

func (lpr _scoper) Scope(
	index int,
	scope map[string]*types.Value,
	scgLoopRange interface{},
	scgLoopVars *types.SCGLoopVars,
	opHandle types.DataHandle,
) (
	map[string]*types.Value,
	error,
) {
	if nil == scgLoopVars {
		return scope, nil
	}

	outboundScope := map[string]*types.Value{}
	for varName, varData := range scope {
		outboundScope[varName] = varData
	}

	if nil != scgLoopVars.Index {
		// assign iteration index to requested inboundScope variable
		indexAsFloat64 := float64(index)
		outboundScope[*scgLoopVars.Index] = &types.Value{
			Number: &indexAsFloat64,
		}
	}

	if nil == scgLoopRange {
		// guard no loopable
		return outboundScope, nil
	}

	var loopable *types.Value
	var err error
	loopable, err = lpr.loopableInterpreter.Interpret(
		scgLoopRange,
		opHandle,
		outboundScope,
	)
	if nil != err {
		return nil, err
	}

	var rawValue interface{}

	if nil != loopable.Array {
		// loopable is array
		rawValue = (*loopable.Array)[index]
	} else {
		// loopable is object
		if index >= len(*loopable.Object) {
			// beyond range
			return outboundScope, nil
		}

		sortedNames := lpr.sortMap(*loopable.Object)
		name := sortedNames[index]
		rawValue = (*loopable.Object)[name]

		if nil != scgLoopVars.Key {
			// only add key to scope if declared
			outboundScope[*scgLoopVars.Key] = &types.Value{String: &name}
		}
	}

	if nil != scgLoopVars.Value {
		var value types.Value
		value, err = lpr.valueInterpreter.Interpret(
			rawValue,
			outboundScope,
			opHandle,
		)
		if nil != err {
			return nil, err
		}

		outboundScope[*scgLoopVars.Value] = &value
	}

	return outboundScope, nil
}
