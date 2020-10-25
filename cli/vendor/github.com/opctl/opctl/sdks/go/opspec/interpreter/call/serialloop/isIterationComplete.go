package serialloop

import (
	"github.com/opctl/opctl/sdks/go/model"
)

// IsIterationComplete tests if an index is within range of the loop
func IsIterationComplete(
	index int,
	dcgSerialLoop *model.SerialLoopCall,
) bool {
	if nil != dcgSerialLoop.Until && *dcgSerialLoop.Until {
		// exit condition provided & met
		return true
	}

	loopRange := dcgSerialLoop.Range
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
