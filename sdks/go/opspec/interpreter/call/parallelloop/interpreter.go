package parallelloop

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		callParallelLoopSpec model.CallParallelLoopSpec,
		scope map[string]*model.Value,
	) (*model.DCGParallelLoopCall, error)
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
	callParallelLoopSpec model.CallParallelLoopSpec,
	scope map[string]*model.Value,
) (*model.DCGParallelLoopCall, error) {
	dcgParallelLoop := model.DCGParallelLoopCall{}

	loopRangeSpec := callParallelLoopSpec.Range
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
