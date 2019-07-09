package parallelloop

import (
	"github.com/opctl/sdk-golang/model"
	"github.com/opctl/sdk-golang/opspec/interpreter/call/predicates"
	"github.com/opctl/sdk-golang/opspec/interpreter/loopable"
)

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

type Interpreter interface {
	Interpret(
		opHandle model.DataHandle,
		scgParallelLoop model.SCGParallelLoop,
		scope map[string]*model.Value,
	) (*model.DCGParallelLoop, error)
}

// NewInterpreter returns an initialized Interpreter instance
func NewInterpreter() Interpreter {
	return &_interpreter{
		loopableInterpreter:   loopable.NewInterpreter(),
		predicatesInterpreter: predicates.NewInterpreter(),
	}
}

type _interpreter struct {
	loopableInterpreter   loopable.Interpreter
	predicatesInterpreter predicates.Interpreter
}

func (itp _interpreter) Interpret(
	opHandle model.DataHandle,
	scgParallelLoop model.SCGParallelLoop,
	scope map[string]*model.Value,
) (*model.DCGParallelLoop, error) {
	dcgParallelLoop := model.DCGParallelLoop{}

	scgLoopRange := scgParallelLoop.Range
	if nil != scgLoopRange {
		dcgLoopRange, err := itp.loopableInterpreter.Interpret(
			scgLoopRange,
			opHandle,
			scope,
		)
		if nil != err {
			return nil, err
		}

		dcgParallelLoop.Range = dcgLoopRange
	}

	return &dcgParallelLoop, nil
}
