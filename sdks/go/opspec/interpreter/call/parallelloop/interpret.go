package parallelloop

import (
	"github.com/opctl/opctl/sdks/go/model"
	"github.com/opctl/opctl/sdks/go/opspec/interpreter/loopable"
)

// Interpret a parallel Loop
func Interpret(
	parallelLoopCallSpec model.ParallelLoopCallSpec,
	scope map[string]*ipld.Node,
) (*model.ParallelLoopCall, error) {
	parallelLoopCall := model.ParallelLoopCall{}

	loopRangeSpec := parallelLoopCallSpec.Range
	if loopRangeSpec != nil {
		dcgLoopRange, err := loopable.Interpret(
			loopRangeSpec,
			scope,
		)
		if err != nil {
			return nil, err
		}

		parallelLoopCall.Range = dcgLoopRange
	}

	return &parallelLoopCall, nil
}
