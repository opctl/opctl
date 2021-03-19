package parallelloop

import (
	"github.com/opctl/opctl/sdks/go/model"
)

// IsIterationComplete tests if an index is within range of the loop
func IsIterationComplete(
	index int,
	parallelLoopCall model.ParallelLoopCall,
) bool {
	loopRange := parallelLoopCall.Range
	if loopRange != nil {
		if loopRange.Array != nil {
			return index >= len(*loopRange.Array)
		}
		if loopRange.Object != nil {
			return index >= len(*loopRange.Object)
		}

		// empty array or object
		return true
	}

	return false
}
