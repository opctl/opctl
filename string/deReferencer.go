package string

import (
	"fmt"
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/interpolater"
)

func newDeReferencer(
	scope map[string]*model.Value,
) interpolater.DeReferencer {
	return _deReferencer{
		data:  data.New(),
		scope: scope,
	}
}

type _deReferencer struct {
	data  data.Data
	scope map[string]*model.Value
}

func (dr _deReferencer) DeReference(
	ref string,
) (string, bool, error) {
	value, ok := dr.scope[ref]
	if !ok {
		// @TODO: replace w/ error once spec supports escapes; for now treat as literal text
		return ref, false, nil
	}

	valueAsString, err := dr.data.CoerceToString(value)
	if nil != err {
		return "", false, fmt.Errorf("Unable to deReference '%v' as string; error was: %v", ref, err.Error())
	}

	return valueAsString, true, nil
}
