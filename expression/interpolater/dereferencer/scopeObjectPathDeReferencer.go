package dereferencer

//go:generate counterfeiter -o ./fakeScopeObjectPathDeReferencer.go --fake-name fakeScopeObjectPathDeReferencer ./ scopeObjectPathDeReferencer

import (
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
	"strings"
)

// scopeObjectPathDeReferencer de references scope object path refs, i.e. refs of the form: $(name.sub.prop)
type scopeObjectPathDeReferencer interface {
	DeReferenceScopeObjectPath(
		ref string,
		scope map[string]*model.Value,
	) (string, bool, error)
}

func newScopeObjectPathDeReferencer() scopeObjectPathDeReferencer {
	return _scopeObjectPathDeReferencer{
		coerce: coerce.New(),
	}
}

type _scopeObjectPathDeReferencer struct {
	coerce coerce.Coerce
}

func (sod _scopeObjectPathDeReferencer) DeReferenceScopeObjectPath(
	ref string,
	scope map[string]*model.Value,
) (string, bool, error) {
	possiblePathIndex := strings.Index(ref, ".")
	if possiblePathIndex > 0 {
		if scopeValue, isPropertyRef := scope[ref[:possiblePathIndex]]; isPropertyRef {
			// scope object ref w/ path
			objectValue, err := sod.coerce.ToObject(scopeValue)
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

			valueAsString, err := sod.coerce.ToString(value)
			if nil != err {
				return "", false, fmt.Errorf("unable to deReference '%v' as string; error was: %v", ref, err.Error())
			}

			return *valueAsString.String, true, nil
		}
	}
	return ref, false, nil
}
