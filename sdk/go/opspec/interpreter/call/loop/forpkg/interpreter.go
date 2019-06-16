package forpkg

import (
	"github.com/opctl/opctl/sdk/go/model"
	"github.com/opctl/opctl/sdk/go/opspec/interpreter/loopable"
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
		loopableInterpreter: loopable.NewInterpreter(),
	}
}

type _interpreter struct {
	loopableInterpreter loopable.Interpreter
}

func (itp _interpreter) Interpret(
	opHandle model.DataHandle,
	scgLoopFor *model.SCGLoopFor,
	scope map[string]*model.Value,
) (*model.DCGLoopFor, error) {
	dcgForEach, err := itp.loopableInterpreter.Interpret(
		scgLoopFor.Each,
		opHandle,
		scope,
	)
	if nil != err {
		return nil, err
	}

	return &model.DCGLoopFor{
		Each:  dcgForEach,
		Value: scgLoopFor.Value,
	}, nil
}
