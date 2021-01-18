package iteration

import (
	"sort"
	"strings"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
)

func refToVariable(ref *string) string {
	return strings.TrimSuffix(strings.TrimPrefix(*ref, "$("), ")")
}

func sortMap(
	m map[string]interface{},
) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}

	sort.Strings(names) //sort keys alphabetically
	return names
}

// Scope scopes loop iteration vars (index, key, value)
func Scope(
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

	var v *model.Value
	var err error
	v, err = loopable.Interpret(
		callSpecLoopRange,
		outboundScope,
	)
	if nil != err {
		return nil, err
	}

	var rawValue interface{}

	if nil != v.Array {
		// loopable is array
		if index >= len(*v.Array) {
			// beyond range
			return outboundScope, nil
		}

		rawValue = (*v.Array)[index]
	} else {
		// loopable is object
		if index >= len(*v.Object) {
			// beyond range
			return outboundScope, nil
		}

		sortedNames := sortMap(*v.Object)
		name := sortedNames[index]
		rawValue = (*v.Object)[name]

		if nil != loopVarsSpec.Key {
			// only add key to scope if declared
			outboundScope[refToVariable(loopVarsSpec.Key)] = &model.Value{String: &name}
		}
	}

	if nil != loopVarsSpec.Value {
		var v model.Value
		v, err = value.Interpret(
			rawValue,
			outboundScope,
		)
		if nil != err {
			return nil, err
		}

		outboundScope[refToVariable(loopVarsSpec.Value)] = &v
	}

	return outboundScope, nil
}
