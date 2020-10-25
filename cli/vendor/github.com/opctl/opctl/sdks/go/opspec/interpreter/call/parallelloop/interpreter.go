package parallelloop

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		parallelLoopCallSpec model.ParallelLoopCallSpec,
		scope map[string]*model.Value,
	) (*model.ParallelLoopCall, error)
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
	parallelLoopCallSpec model.ParallelLoopCallSpec,
	scope map[string]*model.Value,
) (*model.ParallelLoopCall, error) {
	dcgParallelLoop := model.ParallelLoopCall{}

	loopRangeSpec := parallelLoopCallSpec.Range
	if nil != loopRangeSpec {
		dcgLoopRange, err := itp.loopableInterpreter.Interpret(
			loopRangeSpec,
			scope,
		)
		if nil != err {
			return nil, err
		}

		dcgParallelLoop.Range = dcgLoopRange
	}

	return &dcgParallelLoop, nil
}
