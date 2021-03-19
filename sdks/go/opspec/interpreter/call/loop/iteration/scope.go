package iteration

import (
	"sort"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
)

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
	if loopVarsSpec == nil {
		return scope, nil
	}

	outboundScope := map[string]*model.Value{}
	for varName, varData := range scope {
		outboundScope[varName] = varData
	}

	if loopVarsSpec.Index != nil {
		// assign iteration index to requested inboundScope variable
		indexAsFloat64 := float64(index)
		outboundScope[opspec.RefToName(*loopVarsSpec.Index)] = &model.Value{
			Number: &indexAsFloat64,
		}
	}

	if callSpecLoopRange == nil {
		// guard no loopable
		return outboundScope, nil
	}

	var v *model.Value
	var err error
	v, err = loopable.Interpret(
		callSpecLoopRange,
		outboundScope,
	)
	if err != nil {
		return nil, err
	}

	var rawValue interface{}

	if v.Array != nil {
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

		if loopVarsSpec.Key != nil {
			// only add key to scope if declared
			outboundScope[opspec.RefToName(*loopVarsSpec.Key)] = &model.Value{String: &name}
		}
	}

	if loopVarsSpec.Value != nil {
		var v model.Value
		v, err = value.Interpret(
			rawValue,
			outboundScope,
		)
		if err != nil {
			return nil, err
		}

		outboundScope[opspec.RefToName(*loopVarsSpec.Value)] = &v
	}

	return outboundScope, nil
}
