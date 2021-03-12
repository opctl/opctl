package reference

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/opctl/opctl/sdks/go/internal/uniquestring"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/direntry"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/unbracketed"

	"github.com/opctl/opctl/sdks/go/data/coerce"

	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/bracketed"
)

const (
	operator  = '$'
	refOpener = '('
	refCloser = ')'
	RefStart  = string(operator) + string(refOpener)
	RefEnd    = string(refCloser)
)

// Interpret a ref of the form:
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
func Interpret(
	ref string,
	scope map[string]*model.Value,
	opts *model.ReferenceOpts,
) (*model.Value, error) {

	var data *model.Value
	var err error

	ref = strings.TrimSuffix(strings.TrimPrefix(ref, RefStart), RefEnd)
	ref, err = interpolate(
		ref,
		scope,
	)
	if nil != err {
		return nil, err
	}

	data, ref, err = getRootValue(
		ref,
		scope,
		opts,
	)
	if nil != err {
		return nil, err
	}

	_, data, err = rInterpret(
		ref,
		data,
		opts,
	)
	return data, err
}

// interpolate interpolates a ref; refs can be nested at most, one level i.e. $(refOuter$(refInner))
func interpolate(
	ref string,
	scope map[string]*model.Value,
) (string, error) {
	refBuffer := []byte{}
	i := 0

	for i < len(ref) {
		switch {
		case ref[i] == operator:
			nestedRefStartIndex := i + len(RefStart)
			nestedRefEndBracketOffset := strings.Index(ref[nestedRefStartIndex:], RefEnd)
			if nestedRefEndBracketOffset < 0 {
				return "", fmt.Errorf("unable to interpret '%v' as reference; expected ')'", ref[i:])
			}
			nestedRefEndBracketIndex := nestedRefStartIndex + nestedRefEndBracketOffset
			nestedRef := ref[nestedRefStartIndex:nestedRefEndBracketIndex]
			i += len(RefStart) + len(nestedRef) + len(RefEnd)

			var nestedRefRootValue *model.Value
			var err error
			nestedRefRootValue, nestedRef, err = getRootValue(
				nestedRef,
				scope,
				nil,
			)
			if nil != err {
				return "", err
			}

			_, nestedRefValue, err := rInterpret(
				nestedRef,
				nestedRefRootValue,
				nil,
			)
			if nil != err {
				return "", err
			}

			nestedRefValueAsString, err := coerce.ToString(nestedRefValue)
			if nil != err {
				return "", err
			}
			refBuffer = append(refBuffer, *nestedRefValueAsString.String...)

		default:
			refBuffer = append(refBuffer, ref[i])
			i++
		}
	}
	return string(refBuffer), nil
}

func getRootValue(
	ref string,
	scope map[string]*model.Value,
	opts *model.ReferenceOpts,
) (*model.Value, string, error) {
	if strings.HasPrefix(ref, "/") {
		// handle deprecated absolute path reference
		return scope["/"], ref, nil
	}

	if strings.HasPrefix(ref, "./") {
		// handle current directory reference
		return scope["./"], ref[1:], nil
	}

	if strings.HasPrefix(ref, "../") {
		// handle parent directory reference
		return scope["../"], ref[2:], nil
	}

	// scope ref
	var identifier string
	var refRemainder string
	identifier, refRemainder = unbracketed.Parse(ref)

	if value, ok := scope[identifier]; ok {
		return value, refRemainder, nil
	}

	if nil != opts {

		uuid, _ := uniquestring.Construct()
		fsPath := filepath.Join(opts.ScratchDir, uuid)

		switch opts.Type {
		case "Dir":
			os.MkdirAll(fsPath, 0700)
			return &model.Value{Dir: &fsPath}, "", nil
		case "File":
			os.MkdirAll(filepath.Dir(fsPath), 0700)
			os.Create(fsPath)
			return &model.Value{File: &fsPath}, "", nil
		}
	}

	return nil, "", fmt.Errorf("unable to interpret '%v' as reference; '%v' not in scope", ref, identifier)
}

// rInterpret interprets refs of the form:
// .i1
// .i1.i2
// .i1[i2]
// [i1].i2
// [i1][i2]
// i1/p1.ext
func rInterpret(
	ref string,
	data *model.Value,
	opts *model.ReferenceOpts,
) (string, *model.Value, error) {

	if "" == ref {
		return "", data, nil
	}

	switch ref[0] {
	case '[':
		ref, data, err := bracketed.Interpret(ref, data)
		if nil != err {
			return "", nil, err
		}

		return rInterpret(ref, data, opts)
	case '.':
		ref, data, err := unbracketed.Interpret(ref[1:], data)
		if nil != err {
			return "", nil, err
		}

		return rInterpret(ref, data, opts)
	case '/':
		var createType *string
		if nil != opts {
			createType = &opts.Type
		}
		ref, data, err := direntry.Interpret(ref, data, createType)
		if nil != err {
			return "", nil, err
		}

		return rInterpret(ref, data, opts)
	default:
		return "", nil, fmt.Errorf("unable to interpret '%v' as reference; expected '[', '.', or '/'", ref)
	}

}
