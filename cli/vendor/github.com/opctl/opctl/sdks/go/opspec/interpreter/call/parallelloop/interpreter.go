package parallelloop

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		opHandle model.DataHandle,
		scgParallelLoop model.SCGParallelLoopCall,
		scope map[string]*model.Value,
	) (*model.DCGParallelLoopCall, error)
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
	scgParallelLoop model.SCGParallelLoopCall,
	scope map[string]*model.Value,
) (*model.DCGParallelLoopCall, error) {
	dcgParallelLoop := model.DCGParallelLoopCall{}

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
