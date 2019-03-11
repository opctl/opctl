package dereferencer

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ DeReferencer

import (
	"fmt"
	"strings"

	"github.com/opctl/sdk-golang/opspec/interpreter/interpolater/dereferencer/direntry"
	unbracketedIdentifier "github.com/opctl/sdk-golang/opspec/interpreter/interpolater/dereferencer/identifier/unbracketed"

	"github.com/opctl/sdk-golang/data/coerce"

	"github.com/opctl/sdk-golang/model"
	bracketedIdentifier "github.com/opctl/sdk-golang/opspec/interpreter/interpolater/dereferencer/identifier/bracketed"
)

// DeReferencer de references refs of the form:
// /p1.ext
// i1
// .i1
// .i1.i2
// .i1[i2]
// [i1].i2
// [i1][i2]
// i1/p1.ext
// - scope refs: $(name)
// - scope object path refs: $(name.sub.prop)
// - scope file path refs: $(name/sub/file.ext)
// - op file path refs: $(/name/sub/file.ext)
type DeReferencer interface {
	// DeReference returns the de referenced value (if any), whether de referencing occurred, and any err
	DeReference(
		ref string,
		scope map[string]*model.Value,
		opHandle model.DataHandle,
	) (string, bool, error)
}

// New returns a DeReferencer
func New() DeReferencer {
	return _deReferencer{
		opFilePathDeReferencer:            newOpFilePathDeReferencer(),
		coerce:                            coerce.New(),
		dirEntryDeReferencer:              direntry.NewDeReferencer(),
		bracketedIdentifierDeReferencer:   bracketedIdentifier.NewDeReferencer(),
		unbracketedIdentifierDeReferencer: unbracketedIdentifier.NewDeReferencer(),
		unbracketedIdentifierParser:       unbracketedIdentifier.NewParser(),
	}
}

type _deReferencer struct {
	opFilePathDeReferencer
	coerce                            coerce.Coerce
	dirEntryDeReferencer              direntry.DeReferencer
	bracketedIdentifierDeReferencer   bracketedIdentifier.DeReferencer
	unbracketedIdentifierDeReferencer unbracketedIdentifier.DeReferencer
	unbracketedIdentifierParser       unbracketedIdentifier.Parser
}

func (dr _deReferencer) DeReference(
	ref string,
	scope map[string]*model.Value,
	opHandle model.DataHandle,
) (string, bool, error) {

	if strings.HasPrefix(ref, "/") {
		return dr.opFilePathDeReferencer.DeReferenceOpFilePath(ref, scope, opHandle)
	}

	identifier, refRemainder := dr.unbracketedIdentifierParser.Parse(ref)
	data, isInScope := scope[identifier]
	if !isInScope {
		// @TODO: replace w/ error once spec supports escapes; for now treat as literal text
		return ref, false, nil
	}

	if "" != refRemainder {
		var err error
		_, data, err = dr.rDeReference(
			refRemainder,
			data,
		)
		if nil != err {
			return "", true, err
		}
	}

	value, err := dr.coerce.ToString(data)
	if nil != err {
		return "", true, err
	}

	return *value.String, true, nil
}

// rDeReference dereferences refs of the form:
// .i1
// .i1.i2
// .i1[i2]
// [i1].i2
// [i1][i2]
// i1/p1.ext
func (dr _deReferencer) rDeReference(
	ref string,
	data *model.Value,
) (string, *model.Value, error) {

	if "" == ref {
		return "", data, nil
	}

	switch ref[0] {
	case '[':
		ref, data, err := dr.bracketedIdentifierDeReferencer.DeReference(ref, data)
		if nil != err {
			return "", nil, err
		}

		return dr.rDeReference(ref, data)
	case '.':
		ref, data, err := dr.unbracketedIdentifierDeReferencer.DeReference(ref[1:], data)
		if nil != err {
			return "", nil, err
		}

		return dr.rDeReference(ref, data)
	case '/':
		ref, data, err := dr.dirEntryDeReferencer.DeReference(ref, data)
		if nil != err {
			return "", nil, err
		}

		return dr.rDeReference(ref, data)
	default:
		return "", nil, fmt.Errorf("unable to deReference '%v'; expected '[', '.', or '/'", ref)
	}

}
