package string

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/coerce"
	"github.com/opspec-io/sdk-golang/util/interpolater"
)

func newDeReferencer(
	scope map[string]*model.Value,
) interpolater.DeReferencer {
	return _deReferencer{
		coerce: coerce.New(),
		scope:  scope,
	}
}

type _deReferencer struct {
	coerce coerce.Coerce
	scope  map[string]*model.Value
}

func (dr _deReferencer) DeReference(
	ref string,
) (string, bool, error) {
	value, ok := dr.scope[ref]
	if !ok {
		// @TODO: replace w/ error once spec supports escapes; for now treat as literal text
		return ref, false, nil
	}

	stringValue, err := dr.coerce.ToString(value)
	if nil != err {
		return "", false, fmt.Errorf("Unable to deReference '%v' as string; error was: %v", ref, err.Error())
	}

	return stringValue, true, nil
}
