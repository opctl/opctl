package string

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/interpolater"
)

func newDeReferencer(
	scope map[string]*model.Value,
) interpolater.DeReferencer {
	return _deReferencer{
		coercer: newCoercer(),
		scope:   scope,
	}
}

type _deReferencer struct {
	coercer coercer
	scope   map[string]*model.Value
}

func (dr _deReferencer) DeReference(
	ref string,
) (string, error) {
	value, ok := dr.scope[ref]
	if !ok {
		// @TODO: replace w/ error once spec supports escapes; for now treat as literal text
		return ref, nil
	}

	stringValue, err := dr.coercer.Coerce(value)
	if nil != err {
		return "", fmt.Errorf("Unable to deReference '%v' as string; error was: %v", ref, err.Error())
	}

	return stringValue, nil
}
