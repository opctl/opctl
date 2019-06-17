package loop

import (
	"github.com/opctl/sdk-golang/model"
)

// IsIterationComplete tests if an index is within range of the loop
func IsIterationComplete(
	index int,
	loop *model.DCGLoop,
) bool {
	if nil != loop.Until && *loop.Until {
		// exit condition provided & met
		return true
	}

	if nil != loop.For {
		if nil != loop.For.Each.Array {
			return index == len(*loop.For.Each.Array)
		}
		if nil != loop.For.Each.Object {
			return index == len(*loop.For.Each.Object)
		}

		// empty array or object
		return true
	}

	return false
}
