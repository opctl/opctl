package iteration

import (
	"sort"
	"strings"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"

	"github.com/opctl/opctl/sdks/go/model"
)

//counterfeiter:generate -o fakes/scoper.go . Scoper
type Scoper interface {
	// Scope scopes loop iteration vars (index, key, value)
	Scope(
		index int,
		scope map[string]*model.Value,
		callSpecLoopRange interface{},
		loopVarsSpec *model.LoopVarsSpec,
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

func refToVariable(ref *string) string {
	return strings.TrimSuffix(strings.TrimPrefix(*ref, "$("), ")")
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
	callSpecLoopRange interface{},
	loopVarsSpec *model.LoopVarsSpec,
) (
	map[string]*model.Value,
	error,
) {
	if nil == loopVarsSpec {
		return scope, nil
	}

	outboundScope := map[string]*model.Value{}
	for varName, varData := range scope {
		outboundScope[varName] = varData
	}

	if nil != loopVarsSpec.Index {
		// assign iteration index to requested inboundScope variable
		indexAsFloat64 := float64(index)
		outboundScope[refToVariable(loopVarsSpec.Index)] = &model.Value{
			Number: &indexAsFloat64,
		}
	}

	if nil == callSpecLoopRange {
		// guard no loopable
		return outboundScope, nil
	}

	var loopable *model.Value
	var err error
	loopable, err = lpr.loopableInterpreter.Interpret(
		callSpecLoopRange,
		outboundScope,
	)
	if nil != err {
		return nil, err
	}

	var rawValue interface{}

	if nil != loopable.Array {
		// loopable is array
		if index >= len(*loopable.Array) {
			// beyond range
			return outboundScope, nil
		}

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

		if nil != loopVarsSpec.Key {
			// only add key to scope if declared
			outboundScope[refToVariable(loopVarsSpec.Key)] = &model.Value{String: &name}
		}
	}

	if nil != loopVarsSpec.Value {
		var value model.Value
		value, err = lpr.valueInterpreter.Interpret(
			rawValue,
			outboundScope,
		)
		if nil != err {
			return nil, err
		}

		outboundScope[refToVariable(loopVarsSpec.Value)] = &value
	}

	return outboundScope, nil
}
