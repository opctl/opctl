package serialloop

import (
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
	"github.com/opctl/opctl/sdks/go/types"
)

//go:generate counterfeiter -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

type Interpreter interface {
	Interpret(
		opHandle types.DataHandle,
		scgSerialLoop types.SCGSerialLoopCall,
		scope map[string]*types.Value,
	) (*types.DCGSerialLoopCall, error)
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
	scgSerialLoop types.SCGSerialLoopCall,
	scope map[string]*types.Value,
) (*types.DCGSerialLoopCall, error) {
	dcgSerialLoop := types.DCGSerialLoopCall{}

	scgLoopRange := scgSerialLoop.Range
	if nil != scgLoopRange {
		dcgLoopRange, err := itp.loopableInterpreter.Interpret(
			scgLoopRange,
			opHandle,
			scope,
		)
		if nil != err {
			return nil, err
		}

		dcgSerialLoop.Range = dcgLoopRange
	}

	scgLoopUntil := scgSerialLoop.Until
	if nil != scgLoopUntil {
		dcgLoopUntil, err := itp.predicatesInterpreter.Interpret(
			opHandle,
			scgLoopUntil,
			scope,
		)
		if nil != err {
			return nil, err
		}

		dcgSerialLoop.Until = &dcgLoopUntil
	}

	return &dcgSerialLoop, nil
}
