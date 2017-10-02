package interpolater

import (
	"context"
	"fmt"
	"github.com/Jeffail/gabs"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"io/ioutil"
	"strings"
)

// deReferencer de references references
type deReferencer interface {
	// DeReference returns the de referenced value (if any), whether de referencing occurred, and any err
	DeReference(
		ref string,
		scope map[string]*model.Value,
		pkgHandle model.PkgHandle,
	) (string, bool, error)
}

func newDeReferencer() deReferencer {
	return _deReferencer{
		data: data.New(),
	}
}

type _deReferencer struct {
	data data.Data
}

func (dr _deReferencer) DeReference(
	ref string,
	scope map[string]*model.Value,
	pkgHandle model.PkgHandle,
) (string, bool, error) {

	if strings.HasPrefix(ref, "/") {
		// pkg content ref
		contentReadSeekCloser, err := pkgHandle.GetContent(context.TODO(), ref)
		if nil != err {
			return "", false, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
		}
		defer contentReadSeekCloser.Close()

		contentBytes, err := ioutil.ReadAll(contentReadSeekCloser)
		if nil != err {
			return "", false, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
		}
		return string(contentBytes), true, nil
	}

	possibleObjectPathRefIndex := strings.Index(ref, ".")
	if possibleObjectPathRefIndex > 0 {
		objectRef, isObjectPathRef := scope[ref[:possibleObjectPathRefIndex]]
		if isObjectPathRef {
			// scope object ref w/ path
			json, err := gabs.Consume(objectRef.Object)
			if nil != err {
				return "", false, fmt.Errorf("unable to deReference '%v'; error was: %v", ref, err.Error())
			}

			var value *model.Value
			valueAtPath := json.Path(ref[possibleObjectPathRefIndex+1:]).Data()
			switch valueAtPath := valueAtPath.(type) {
			case float64:
				value = &model.Value{Number: &valueAtPath}
			case map[string]interface{}:
				value = &model.Value{Object: valueAtPath}
			case string:
				return valueAtPath, true, nil
			case []interface{}:
				value = &model.Value{Array: valueAtPath}
			}

			valueAsString, err := dr.data.CoerceToString(value)
			if nil != err {
				return "", false, fmt.Errorf("unable to deReference '%v' as string; error was: %v", ref, err.Error())
			}

			return *valueAsString.String, true, nil
		}
	}

	// scope ref
	value, ok := scope[ref]
	if !ok {
		// @TODO: replace w/ error once spec supports escapes; for now treat as literal text
		return ref, false, nil
	}

	valueAsString, err := dr.data.CoerceToString(value)
	if nil != err {
		return "", false, fmt.Errorf("unable to deReference '%v' as string; error was: %v", ref, err.Error())
	}

	return *valueAsString.String, true, nil
}
