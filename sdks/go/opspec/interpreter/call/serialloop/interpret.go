package serialloop

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/call/predicates"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
)

//Interpret a serial loop
func Interpret(
	serialLoopCallSpec model.SerialLoopCallSpec,
	scope map[string]*model.Value,
) (*model.SerialLoopCall, error) {
	dcgSerialLoop := model.SerialLoopCall{}

	loopRangeSpec := serialLoopCallSpec.Range
	if nil != loopRangeSpec {
		dcgLoopRange, err := loopable.Interpret(
			loopRangeSpec,
			scope,
		)
		if nil != err {
			return nil, err
		}

		dcgSerialLoop.Range = dcgLoopRange
	}

	callSpecLoopUntil := serialLoopCallSpec.Until
	if nil != callSpecLoopUntil {
		dcgLoopUntil, err := predicates.Interpret(
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
