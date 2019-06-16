package reference

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

import (
	"fmt"
	"strings"

	"github.com/opctl/opctl/sdk/go/opspec/interpreter/reference/direntry"
	unbracketedIdentifier "github.com/opctl/opctl/sdk/go/opspec/interpreter/reference/identifier/unbracketed"

	"github.com/opctl/opctl/sdk/go/data/coerce"

	"github.com/opctl/opctl/sdk/go/model"
	bracketedIdentifier "github.com/opctl/opctl/sdk/go/opspec/interpreter/reference/identifier/bracketed"
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
		scope map[string]*model.Value,
		opHandle model.DataHandle,
	) (*model.Value, error)
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

func (dr _interpreter) Interpret(
	ref string,
	scope map[string]*model.Value,
	opHandle model.DataHandle,
) (*model.Value, error) {

	var data *model.Value
	var refRemainder string
	var err error

	ref = strings.TrimSuffix(strings.TrimPrefix(ref, "$("), ")")

	if strings.HasPrefix(ref, "/") {
		// op path ref
		data = &model.Value{Dir: opHandle.Path()}
		refRemainder = ref
	} else {
		// scope ref
		var identifier string
		identifier, refRemainder = dr.unbracketedIdentifierParser.Parse(ref)

		var isInScope bool
		data, isInScope = scope[identifier]
		if !isInScope {
			return nil, fmt.Errorf("unable to interpret '%v' as reference; '%v' not in scope", identifier, identifier)
		}
	}

	_, data, err = dr.rInterpret(
		refRemainder,
		data,
	)
	return data, err
}

// rInterpret interprets refs of the form:
// .i1
// .i1.i2
// .i1[i2]
// [i1].i2
// [i1][i2]
// i1/p1.ext
func (dr _interpreter) rInterpret(
	ref string,
	data *model.Value,
) (string, *model.Value, error) {

	if "" == ref {
		return "", data, nil
	}

	switch ref[0] {
	case '[':
		ref, data, err := dr.bracketedIdentifierInterpreter.Interpret(ref, data)
		if nil != err {
			return "", nil, err
		}

		return dr.rInterpret(ref, data)
	case '.':
		ref, data, err := dr.unbracketedIdentifierInterpreter.Interpret(ref[1:], data)
		if nil != err {
			return "", nil, err
		}

		return dr.rInterpret(ref, data)
	case '/':
		ref, data, err := dr.dirEntryInterpreter.Interpret(ref, data)
		if nil != err {
			return "", nil, err
		}

		return dr.rInterpret(ref, data)
	default:
		return "", nil, fmt.Errorf("unable to interpret '%v'; expected '[', '.', or '/'", ref)
	}

}
