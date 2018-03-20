package dereferencer

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ DeReferencer

import (
	"github.com/opspec-io/sdk-golang/data/coerce"
	"github.com/opspec-io/sdk-golang/model"
)

// DeReferencer de references:
// - scope refs: $(name)
// - scope object path refs: $(name.sub.prop)
// - scope file path refs: $(name/sub/file.ext)
// - pkg file path refs: $(/name/sub/file.ext)
type DeReferencer interface {
	// DeReference returns the de referenced value (if any), whether de referencing occurred, and any err
	DeReference(
		ref string,
		scope map[string]*model.Value,
		opDirHandle model.DataHandle,
	) (string, bool, error)
}

// New returns a DeReferencer
func New() DeReferencer {
	return _deReferencer{
		coerce:                      coerce.New(),
		pkgFilePathDeReferencer:     newPkgFilePathDeReferencer(),
		scopeDeReferencer:           newScopeDeReferencer(),
		scopeFilePathDeReferencer:   newScopeFilePathDeReferencer(),
		scopeObjectPathDeReferencer: newScopeObjectPathDeReferencer(),
	}
}

type _deReferencer struct {
	coerce coerce.Coerce
	pkgFilePathDeReferencer
	scopeDeReferencer
	scopeFilePathDeReferencer
	scopeObjectPathDeReferencer
}

func (dr _deReferencer) DeReference(
	ref string,
	scope map[string]*model.Value,
	opDirHandle model.DataHandle,
) (string, bool, error) {

	var (
		value string
		err   error
	)

	var isPkgFilePathRef bool
	if value, isPkgFilePathRef, err = dr.pkgFilePathDeReferencer.DeReferencePkgFilePath(ref, scope, opDirHandle); isPkgFilePathRef {
		return value, isPkgFilePathRef, err
	}

	var isScopeRef bool
	if value, isScopeRef, err = dr.scopeDeReferencer.DeReferenceScope(ref, scope); isScopeRef {
		return value, isScopeRef, err
	}

	var isScopeFilePathRef bool
	if value, isScopeFilePathRef, err = dr.scopeFilePathDeReferencer.DeReferenceScopeFilePath(ref, scope); isScopeFilePathRef {
		return value, isScopeFilePathRef, err
	}

	var isScopeObjectPathRef bool
	if value, isScopeObjectPathRef, err = dr.scopeObjectPathDeReferencer.DeReferenceScopeObjectPath(ref, scope); isScopeObjectPathRef {
		return value, isScopeObjectPathRef, err
	}

	// @TODO: replace w/ error once spec supports escapes; for now treat as literal text
	return ref, false, nil
}
