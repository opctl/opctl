package reference

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"strings"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/direntry"
	unbracketedIdentifier "github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/unbracketed"

	"github.com/opctl/opctl/sdks/go/data/coerce"

	bracketedIdentifier "github.com/opctl/opctl/sdks/go/opspec/interpreter/reference/identifier/bracketed"
	"github.com/opctl/opctl/sdks/go/types"
)

const (
	operator  = '$'
	refOpener = '('
	refCloser = ')'
	RefStart  = string(operator) + string(refOpener)
	RefEnd    = string(refCloser)
)

// Interpreter interprets refs of the form:
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
type Interpreter interface {
	// Interpret returns the interpreted value (if any), and any err
	Interpret(
		ref string,
		scope map[string]*types.Value,
		opHandle types.DataHandle,
	) (*types.Value, error)
}

// NewInterpreter returns a Interpreter
func NewInterpreter() Interpreter {
	return _interpreter{
		coerce:                           coerce.New(),
		dirEntryInterpreter:              direntry.NewInterpreter(),
		bracketedIdentifierInterpreter:   bracketedIdentifier.NewInterpreter(),
		unbracketedIdentifierInterpreter: unbracketedIdentifier.NewInterpreter(),
		unbracketedIdentifierParser:      unbracketedIdentifier.NewParser(),
	}
}

type _interpreter struct {
	coerce                           coerce.Coerce
	dirEntryInterpreter              direntry.Interpreter
	bracketedIdentifierInterpreter   bracketedIdentifier.Interpreter
	unbracketedIdentifierInterpreter unbracketedIdentifier.Interpreter
	unbracketedIdentifierParser      unbracketedIdentifier.Parser
}

func (itp _interpreter) Interpret(
	ref string,
	scope map[string]*types.Value,
	opHandle types.DataHandle,
) (*types.Value, error) {

	var data *types.Value
	var err error

	ref = strings.TrimSuffix(strings.TrimPrefix(ref, RefStart), RefEnd)
	ref, err = itp.interpolate(
		ref,
		scope,
		opHandle,
	)
	if nil != err {
		return nil, err
	}

	data, ref, err = itp.getRootValue(
		ref,
		scope,
		opHandle,
	)
	if nil != err {
		return nil, err
	}

	_, data, err = itp.rInterpret(
		ref,
		data,
	)
	return data, err
}

// interpolate interpolates a ref; refs can be nested at most, one level i.e. $(refOuter$(refInner))
func (itp _interpreter) interpolate(
	ref string,
	scope map[string]*types.Value,
	opHandle types.DataHandle,
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

			var nestedRefRootValue *types.Value
			var err error
			nestedRefRootValue, nestedRef, err = itp.getRootValue(
				nestedRef,
				scope,
				opHandle,
			)
			if nil != err {
				return "", err
			}

			_, nestedRefValue, err := itp.rInterpret(
				nestedRef,
				nestedRefRootValue,
			)
			if nil != err {
				return "", err
			}

			nestedRefValueAsString, err := itp.coerce.ToString(nestedRefValue)
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

func (itp _interpreter) getRootValue(
	ref string,
	scope map[string]*types.Value,
	opHandle types.DataHandle,
) (*types.Value, string, error) {
	if strings.HasPrefix(ref, "/") {
		// op path ref
		return &types.Value{Dir: opHandle.Path()}, ref, nil
	}

	// scope ref
	var identifier string
	var refRemainder string
	identifier, refRemainder = itp.unbracketedIdentifierParser.Parse(ref)

	if value, ok := scope[identifier]; ok {
		return value, refRemainder, nil
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
func (itp _interpreter) rInterpret(
	ref string,
	data *types.Value,
) (string, *types.Value, error) {

	if "" == ref {
		return "", data, nil
	}

	switch ref[0] {
	case '[':
		ref, data, err := itp.bracketedIdentifierInterpreter.Interpret(ref, data)
		if nil != err {
			return "", nil, err
		}

		return itp.rInterpret(ref, data)
	case '.':
		ref, data, err := itp.unbracketedIdentifierInterpreter.Interpret(ref[1:], data)
		if nil != err {
			return "", nil, err
		}

		return itp.rInterpret(ref, data)
	case '/':
		ref, data, err := itp.dirEntryInterpreter.Interpret(ref, data)
		if nil != err {
			return "", nil, err
		}

		return itp.rInterpret(ref, data)
	default:
		return "", nil, fmt.Errorf("unable to interpret '%v' as reference; expected '[', '.', or '/'", ref)
	}

}
