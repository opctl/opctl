package parallelloop

import (
	"github.com/opctl/sdk-golang/model"
)

// IsIterationComplete tests if an index is within range of the loop
func IsIterationComplete(
	index int,
	dcgParallelLoop model.DCGParallelLoop,
) bool {
	loopRange := dcgParallelLoop.Range
	if nil != loopRange {
		if nil != loopRange.Array {
			return index == len(*loopRange.Array)
		}
		if nil != loopRange.Object {
			return index == len(*loopRange.Object)
		}

		// empty array or object
		return true
	}

	return false
}
