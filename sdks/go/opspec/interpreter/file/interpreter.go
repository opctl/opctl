package file

import (
	"fmt"
	"strings"

	"github.com/opctl/opctl/sdks/go/opspec/interpreter/interpolater"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/reference"

	"github.com/opctl/opctl/sdks/go/data/coerce"
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/value"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	// Interpret interprets an expression to a file value.
	// Expression must be a type supported by coerce.ToFile
	// scratchDir will be used as the containing dir if file creation necessary
	//
	// Examples of valid file expressions:
	// scope ref: $(scope-ref)
	// scope ref w/ path: $(scope-ref/file.txt)
	// pkg fs ref: $(/pkg-fs-ref)
	// pkg fs ref w/ path: $(/pkg-fs-ref/file.txt)
	Interpret(
		scope map[string]*model.Value,
		expression interface{},
		opHandle model.DataHandle,
		scratchDir string,
	) (*model.Value, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return _interpreter{
		coerce:               coerce.New(),
		referenceInterpreter: reference.NewInterpreter(),
		valueInterpreter:     value.NewInterpreter(),
	}
}

type _interpreter struct {
	coerce               coerce.Coerce
	referenceInterpreter reference.Interpreter
	valueInterpreter     value.Interpreter
}

func (itp _interpreter) Interpret(
	scope map[string]*model.Value,
	expression interface{},
	opHandle model.DataHandle,
	scratchDir string,
) (*model.Value, error) {
	expressionAsString, expressionIsString := expression.(string)

	// @TODO: this incorrectly treats $(inScope)$(inScope) as ref
	if expressionIsString && strings.HasPrefix(expressionAsString, interpolater.RefStart) && strings.HasSuffix(expressionAsString, interpolater.RefEnd) {
		value, err := itp.referenceInterpreter.Interpret(
			expressionAsString,
			scope,
			opHandle,
		)
		if nil != err {
			return nil, fmt.Errorf("unable to interpret %+v to file; error was %v", expression, err)
		}
		return itp.coerce.ToFile(value, scratchDir)
	}

	value, err := itp.valueInterpreter.Interpret(
		expression,
		scope,
		opHandle,
	)
	if nil != err {
		return nil, fmt.Errorf("unable to interpret %+v to file; error was %v", expression, err)
	}

	return itp.coerce.ToFile(&value, scratchDir)
}
