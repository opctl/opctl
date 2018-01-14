package dereferencer

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ DeReferencer

import (
	"github.com/opspec-io/sdk-golang/data"
	"github.com/opspec-io/sdk-golang/model"
)

// DeReferencer de references:
// - scope refs: $(name)
// - scope property refs: $(name.sub.prop)
// - scope file refs: $(name/sub/file.ext)
// - pkg file refs: $(/name/sub/file.ext)
type DeReferencer interface {
	// DeReference returns the de referenced value (if any), whether de referencing occurred, and any err
	DeReference(
		ref string,
		scope map[string]*model.Value,
		pkgHandle model.PkgHandle,
	) (string, bool, error)
}

// New returns a DeReferencer
func New() DeReferencer {
	return _deReferencer{
		data:                      data.New(),
		pkgFileDeReferencer:       newPkgFileDeReferencer(),
		scopeDeReferencer:         newScopeDeReferencer(),
		scopeFileDeReferencer:     newScopeFileDeReferencer(),
		scopePropertyDeReferencer: newScopePropertyDeReferencer(),
	}
}

type _deReferencer struct {
	data data.Data
	pkgFileDeReferencer
	scopeDeReferencer
	scopeFileDeReferencer
	scopePropertyDeReferencer
}

func (dr _deReferencer) DeReference(
	ref string,
	scope map[string]*model.Value,
	pkgHandle model.PkgHandle,
) (string, bool, error) {

	var (
		value string
		err   error
	)

	var isPkgFileRef bool
	if value, isPkgFileRef, err = dr.pkgFileDeReferencer.DeReferencePkgFile(ref, scope, pkgHandle); isPkgFileRef {
		return value, isPkgFileRef, err
	}

	var isScopeRef bool
	if value, isScopeRef, err = dr.scopeDeReferencer.DeReferenceScope(ref, scope); isScopeRef {
		return value, isScopeRef, err
	}

	var isScopeFileRef bool
	if value, isScopeFileRef, err = dr.scopeFileDeReferencer.DeReferenceScopeFile(ref, scope); isScopeFileRef {
		return value, isScopeFileRef, err
	}

	var isScopePropertyRef bool
	if value, isScopePropertyRef, err = dr.scopePropertyDeReferencer.DeReferenceScopeProperty(ref, scope); isScopePropertyRef {
		return value, isScopePropertyRef, err
	}

	// @TODO: replace w/ error once spec supports escapes; for now treat as literal text
	return ref, false, nil
}
