package dereferencer

//go:generate counterfeiter -o ./fakeScopeDeReferencer.go --fake-name fakeScopeDeReferencer ./ scopeDeReferencer

import (
	"fmt"

	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
)

// scopeDeReferencer de references scope refs, i.e. refs of the form: $(name)
type scopeDeReferencer interface {
	DeReferenceScope(
		ref string,
		scope map[string]*model.Value,
	) (string, bool, error)
}

func newScopeDeReferencer() scopeDeReferencer {
	return _scopeDeReferencer{
		data: data.New(),
	}
}

type _scopeDeReferencer struct {
	data data.Data
}

func (sd _scopeDeReferencer) DeReferenceScope(
	ref string,
	scope map[string]*model.Value,
) (string, bool, error) {
	if value, isScopeRef := scope[ref]; isScopeRef {
		valueAsString, err := sd.data.CoerceToString(value)
		if nil != err {
			return "", false, fmt.Errorf("unable to deReference '%v' as string; error was: %v", ref, err.Error())
		}

		return *valueAsString.String, true, nil
	}

	return ref, false, nil
}
