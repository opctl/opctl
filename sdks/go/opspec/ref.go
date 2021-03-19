package opspec

import (
	"fmt"
	"strings"
)

// RefToName converts a variable reference to the name of the variable
func RefToName(ref string) string {
	return strings.TrimSuffix(strings.TrimPrefix(ref, "$("), ")")
}

// NameToRef converts a variable name to the reference form in an opspec
func NameToRef(name string) string {
	return fmt.Sprintf("$(%s)", name)
}
