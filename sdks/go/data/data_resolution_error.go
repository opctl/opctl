package data

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// ErrDataResolution aggregates errors from all attempts to resolve data into
// one more user-friendly one
type ErrDataResolution struct {
	dataRef string
	errs    []error
}

func (err ErrDataResolution) Error() string {
	messageBuffer := bytes.NewBufferString(fmt.Sprintf("unable to resolve op \"%s\":", err.dataRef))

	for _, err := range err.errs {
		parts := strings.Split(err.Error(), "\n")
		if len(parts) > 1 {
			for i, part := range parts {
				prefix := " "
				if i == 0 {
					prefix = "-"
				}
				messageBuffer.WriteString(fmt.Sprintf("\n%s %s", prefix, part))
			}
		} else {
			messageBuffer.WriteString(fmt.Sprintf("\n- %v", err))
		}
	}

	return messageBuffer.String()
}

// Is reports whether any of the internal errors matches the target
func (err ErrDataResolution) Is(target error) bool {
	for _, err := range err.errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
