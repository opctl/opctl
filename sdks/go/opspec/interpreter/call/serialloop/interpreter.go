package serialloop

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -o ./fakeInterpreter.go --fake-name FakeInterpreter ./ Interpreter

type Interpreter interface {
	Interpret(
		opHandle model.DataHandle,
		scgSerialLoop model.SCGSerialLoopCall,
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
	opHandle model.DataHandle,
	scgSerialLoop model.SCGSerialLoopCall,
	scope map[string]*model.Value,
) (*model.DCGSerialLoopCall, error) {
	dcgSerialLoop := model.DCGSerialLoopCall{}

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
