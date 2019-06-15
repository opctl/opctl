package iteration

//go:generate counterfeiter -o ./fakeScoper.go --fake-name FakeScoper ./ Scoper

import (
	"sort"

	"github.com/opctl/sdk-golang/opspec/interpreter/value"

	"github.com/opctl/sdk-golang/opspec/interpreter/loopable"

	"github.com/opctl/sdk-golang/model"
)

type Scoper interface {
	// Scope scopes loop iteration vars (index, key, value)
	Scope(
		index int,
		scope map[string]*model.Value,
		scgLoop *model.SCGLoop,
		opHandle model.DataHandle,
	) (
		map[string]*model.Value,
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
	scope map[string]*model.Value,
	scgLoop *model.SCGLoop,
	opHandle model.DataHandle,
) (
	map[string]*model.Value,
	error,
) {
	outboundScope := map[string]*model.Value{}
	for varName, varData := range scope {
		outboundScope[varName] = varData
	}

	if nil != scgLoop.Index {
		// assign iteration index to requested inboundScope variable
		indexAsFloat64 := float64(index)
		outboundScope[*scgLoop.Index] = &model.Value{
			Number: &indexAsFloat64,
		}
	}

	var loopable *model.Value
	if nil != scgLoop.For && nil != scgLoop.For.Each {
		var err error
		loopable, err = lpr.loopableInterpreter.Interpret(
			scgLoop.For.Each,
			opHandle,
			outboundScope,
		)
		if nil != err {
			return nil, err
		}

		var rawValue interface{}

		if nil != loopable.Array {
			rawValue = (*loopable.Array)[index]
		} else {
			sortedNames := lpr.sortMap(*loopable.Object)
			name := sortedNames[index]
			rawValue = (*loopable.Object)[name]

			if nil != scgLoop.For.Key {
				// only add key to scope if declared
				outboundScope[*scgLoop.For.Key] = &model.Value{String: &name}
			}
		}

		if nil != scgLoop.For.Value {
			var value model.Value
			value, err = lpr.valueInterpreter.Interpret(
				rawValue,
				outboundScope,
				opHandle,
			)
			if nil != err {
				return nil, err
			}

			outboundScope[*scgLoop.For.Value] = &value
		}
	}

	return outboundScope, nil
}
