package parallelloop

import (
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
	"github.com/opctl/opctl/sdks/go/types"
)

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

type Interpreter interface {
	Interpret(
		opHandle types.DataHandle,
		scgParallelLoop types.SCGParallelLoopCall,
		scope map[string]*types.Value,
	) (*types.DCGParallelLoopCall, error)
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
	opHandle types.DataHandle,
	scgParallelLoop types.SCGParallelLoopCall,
	scope map[string]*types.Value,
) (*types.DCGParallelLoopCall, error) {
	dcgParallelLoop := types.DCGParallelLoopCall{}

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
