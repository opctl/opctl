package uniquestring

import (
	"github.com/satori/go.uuid"
	"strings"
)

// Construct returns a globally unique string
func Construct() (string, error) {
	uuid := uuid.NewV4()

	return strings.Replace(
		uuid.String(),
		"-",
		"",
		-1,
	), nil
}
