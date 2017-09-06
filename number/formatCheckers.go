package number

import (
	"math"
)

// IntegerFormatChecker is a gojsonschema.FormatChecker for integers
type IntegerFormatChecker struct{}

// Implement gojsonschema.FormatChecker interface
func (f IntegerFormatChecker) IsFormat(input interface{}) bool {

	asFloat64, ok := input.(float64)
	if ok == false {
		return false
	}

	return math.Ceil(asFloat64) == asFloat64
}
