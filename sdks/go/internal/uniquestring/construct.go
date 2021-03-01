package uniquestring

import (
	"github.com/satori/go.uuid"
	"strings"
)

// Construct returns a globally unique string
func Construct() (string, error) {
	uuid, err := uuid.NewV4()
	if nil != err {
		return "", err
	}

	return strings.Replace(
		uuid.String(),
		"-",
		"",
		-1,
	), nil
}
