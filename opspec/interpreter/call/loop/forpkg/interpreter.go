package forpkg

import (
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/array"
)

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

type Interpreter interface {
	Interpret(
		opHandle model.DataHandle,
		scgLoopFor *model.SCGLoopFor,
		scope map[string]*model.Value,
	) (*model.DCGLoopFor, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return &_interpreter{
		arrayInterpreter: array.NewInterpreter(),
	}
}

type _interpreter struct {
	arrayInterpreter array.Interpreter
}

func (itp _interpreter) Interpret(
	opHandle model.DataHandle,
	scgLoopFor *model.SCGLoopFor,
	scope map[string]*model.Value,
) (*model.DCGLoopFor, error) {
	// @TODO: consider an iterableInterpreter to iterate over anything iterable; not only arrays
	dcgForEach, err := itp.arrayInterpreter.Interpret(
		scope,
		scgLoopFor.Each,
		opHandle,
	)
	if nil != err {
		return nil, err
	}

	return &model.DCGLoopFor{
		Each:  dcgForEach,
		Value: scgLoopFor.Value,
	}, nil
}
