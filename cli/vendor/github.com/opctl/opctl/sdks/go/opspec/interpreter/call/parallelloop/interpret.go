package parallelloop

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
)

//Interpret a parallel Loop
func Interpret(
	parallelLoopCallSpec model.ParallelLoopCallSpec,
	scope map[string]*model.Value,
) (*model.ParallelLoopCall, error) {
	parallelLoopCall := model.ParallelLoopCall{}

	loopRangeSpec := parallelLoopCallSpec.Range
	if nil != loopRangeSpec {
		dcgLoopRange, err := loopable.Interpret(
			loopRangeSpec,
			scope,
		)
		if nil != err {
			return nil, err
		}

		parallelLoopCall.Range = dcgLoopRange
	}

	return &parallelLoopCall, nil
}
