package serialloop

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
)

//counterfeiter:generate -o fakes/interpreter.go . Interpreter
type Interpreter interface {
	Interpret(
		callSerialLoopSpec model.CallSerialLoopSpec,
		scope map[string]*model.Value,
	) (*model.DCGSerialLoopCall, error)
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
	callSerialLoopSpec model.CallSerialLoopSpec,
	scope map[string]*model.Value,
) (*model.DCGSerialLoopCall, error) {
	dcgSerialLoop := model.DCGSerialLoopCall{}

	loopRangeSpec := callSerialLoopSpec.Range
	if nil != loopRangeSpec {
		dcgLoopRange, err := itp.loopableInterpreter.Interpret(
			loopRangeSpec,
			scope,
		)
		if nil != err {
			return nil, err
		}

		dcgSerialLoop.Range = dcgLoopRange
	}

	callSpecLoopUntil := callSerialLoopSpec.Until
	if nil != callSpecLoopUntil {
		dcgLoopUntil, err := itp.predicatesInterpreter.Interpret(
			callSpecLoopUntil,
			scope,
		)
		if nil != err {
			return nil, err
		}

		dcgSerialLoop.Until = &dcgLoopUntil
	}

	return &dcgSerialLoop, nil
}
