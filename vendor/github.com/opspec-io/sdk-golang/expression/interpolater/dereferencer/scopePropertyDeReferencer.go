package dereferencer

//go:generate counterfeiter -o ./fakeScopePropertyDeReferencer.go --fake-name fakeScopePropertyDeReferencer ./ scopePropertyDeReferencer

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"strings"
)

// scopePropertyDeReferencer de references scope property refs, i.e. refs of the form: $(name.sub.prop)
type scopePropertyDeReferencer interface {
	DeReferenceScopeProperty(
		ref string,
		scope map[string]*model.Value,
	) (string, bool, error)
}

func newScopePropertyDeReferencer() scopePropertyDeReferencer {
	return _scopePropertyDeReferencer{
		data: data.New(),
	}
}

type _scopePropertyDeReferencer struct {
	data data.Data
}

func (spd _scopePropertyDeReferencer) DeReferenceScopeProperty(
	ref string,
	scope map[string]*model.Value,
) (string, bool, error) {
	possiblePathIndex := strings.Index(ref, ".")
	if possiblePathIndex > 0 {
		if scopeValue, isPropertyRef := scope[ref[:possiblePathIndex]]; isPropertyRef {
			// scope object ref w/ path
			objectValue, err := spd.data.CoerceToObject(scopeValue)
			if nil != err {
				return "", false, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
			}

			json, err := gabs.Consume(objectValue.Object)
			if nil != err {
				return "", false, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
			}

			propertyPath := ref[possiblePathIndex+1:]
			if !json.ExistsP(propertyPath) {
				return "", false, fmt.Errorf("unable to deReference '%v'; path doesn't exist", ref)
			}

			var value *model.Value
			switch valueAtPath := json.Path(propertyPath).Data().(type) {
			case float64:
				value = &model.Value{Number: &valueAtPath}
			case map[string]interface{}:
				value = &model.Value{Object: valueAtPath}
			case string:
				return valueAtPath, true, nil
			case []interface{}:
				value = &model.Value{Array: valueAtPath}
			}

			valueAsString, err := spd.data.CoerceToString(value)
			if nil != err {
				return "", false, fmt.Errorf("unable to deReference '%v' as string; error was: %v", ref, err.Error())
			}

			return *valueAsString.String, true, nil
		}
	}
	return ref, false, nil
}
