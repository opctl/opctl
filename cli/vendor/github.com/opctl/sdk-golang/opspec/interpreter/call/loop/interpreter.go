package loop

import (
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/loop/forpkg"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/predicates"
)

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

type Interpreter interface {
	Interpret(
		opHandle model.DataHandle,
		scgLoop *model.SCGLoop,
		scope map[string]*model.Value,
	) (*model.DCGLoop, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return &_interpreter{
		forInterpreter:        forpkg.NewInterpreter(),
		predicatesInterpreter: predicates.NewInterpreter(),
	}
}

type _interpreter struct {
	forInterpreter        forpkg.Interpreter
	predicatesInterpreter predicates.Interpreter
}

func (itp _interpreter) Interpret(
	opHandle model.DataHandle,
	scgLoop *model.SCGLoop,
	scope map[string]*model.Value,
) (*model.DCGLoop, error) {
	dcgLoop := model.DCGLoop{
		Index: scgLoop.Index,
	}

	if nil != scgLoop.For {
		dcgLoopFor, err := itp.forInterpreter.Interpret(
			opHandle,
			scgLoop.For,
			scope,
		)
		if nil != err {
			return nil, err
		}

		dcgLoop.For = dcgLoopFor
	}

	if nil != scgLoop.Until {
		dcgLoopUntil, err := itp.predicatesInterpreter.Interpret(
			opHandle,
			scgLoop.Until,
			scope,
		)
		if nil != err {
			return nil, err
		}

		dcgLoop.Until = &dcgLoopUntil
	}

	return &dcgLoop, nil
}
