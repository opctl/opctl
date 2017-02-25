package interpolate

import (
	"strconv"
)

// interpolates a string w/ a number according to the opspec spec
func NumberValue(s string, varName string, varValue float64) string {
	return StringValue(s, varName, strconv.FormatFloat(varValue, 'f', -1, 64))
}
