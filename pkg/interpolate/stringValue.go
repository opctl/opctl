package interpolate

import (
	"fmt"
	"strings"
)

// interpolates a string w/ a string according to the opspec spec
func StringValue(s string, varName string, varValue string) string {
	return strings.Replace(s, fmt.Sprintf(`$(%v)`, varName), varValue, -1)
}
