package interpolater

import (
	"fmt"
	"strings"
)

func Interpolate(input string, varName string, varValue string) string {
	return strings.Replace(input, fmt.Sprintf(`$(%v)`, varName), varValue, -1)
}
