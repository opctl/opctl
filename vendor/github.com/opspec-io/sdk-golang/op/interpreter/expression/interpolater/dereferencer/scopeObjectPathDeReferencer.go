package dereferencer

//go:generate counterfeiter -o ./fakeScopeObjectPathDeReferencer.go --fake-name fakeScopeObjectPathDeReferencer ./ scopeObjectPathDeReferencer

import (
	"fmt"
	"strings"

	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
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
	refParts := strings.SplitN(ref, ".", 2)
	if len(refParts) > 1 {
		if scopeValue, isPropertyRef := scope[refParts[0]]; isPropertyRef {
			// scope object ref w/ path
			value, err := sod.walk(refParts[1], scopeValue.Object)
			if nil != err {
				return "", true, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
			}

			valueAsString, err := sod.coerce.ToString(value)
			if nil != err {
				return "", true, fmt.Errorf("unable to deReference '%v' as string; error was: %v", ref, err.Error())
			}

			return *valueAsString.String, true, nil
		}
	}
	return ref, false, nil
}

func (sod _scopeObjectPathDeReferencer) walk(
	ref string,
	data interface{},
) (*model.Value, error) {
	refParts := strings.SplitN(ref, ".", 2)

	switch value := data.(type) {
	case bool:
		if len(refParts) == 1 {
			return &model.Value{Boolean: &value}, nil
		}
		// if unprocessed refParts the ref is invalid
		return nil, fmt.Errorf("unable to deReference '%v'; path doesn't exist", ref)
	case float64:
		if len(refParts) == 1 {
			return &model.Value{Number: &value}, nil
		}
		// if unprocessed refParts the ref is invalid
		return nil, fmt.Errorf("unable to deReference '%v'; path doesn't exist", ref)
	case int:
		// reprocess as float64
		return sod.walk(ref, float64(value))
	case string:
		if len(refParts) == 1 {
			return &model.Value{String: &value}, nil
		}
		// if unprocessed refParts the ref is invalid
		return nil, fmt.Errorf("unable to deReference '%v'; path doesn't exist", ref)
	case map[string]interface{}:
		if "" == ref {
			// no path part
			return &model.Value{Object: value}, nil
		}

		var nextRef string
		if len(refParts) == 1 {
			// last path part
			nextRef = ""
		} else {
			// not last path part
			nextRef = refParts[1]
		}

		if data, isData := value[refParts[0]]; isData {
			return sod.walk(nextRef, data)
		}

		// if no data at path the ref is invalid
		return nil, fmt.Errorf("unable to deReference '%v'; path doesn't exist", ref)
	case []interface{}:
		if len(refParts) == 1 {
			return &model.Value{Array: value}, nil
		}
		return sod.walk(refParts[1], value)
	default:
		return nil, fmt.Errorf("unable to deReference '%v'; '%+v' unexpected type", ref, data)
	}
}
