package interpolate

import (
	"fmt"
	"strings"
)

// interpolates a string according to the opspec spec
func Interpolate(input string, varName string, varValue string) string {
	return strings.Replace(input, fmt.Sprintf(`$(%v)`, varName), varValue, -1)
}
