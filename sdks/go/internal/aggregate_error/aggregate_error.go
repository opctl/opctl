package errors

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// ErrAggregate aggregates multiple errors into one, in a tree format
// instead of a linked list like `.Wrap` does.
//
// They should always be wrapped with `errors.Wrap(err, "label")` to give them
// context.
type ErrAggregate struct {
	errs []error
}

func (e *ErrAggregate) AddError(err error) {
	e.errs = append(e.errs, err)
}

func (err ErrAggregate) Error() string {
	messageBuffer := bytes.NewBufferString("")

	for _, err := range err.errs {
		parts := strings.Split(err.Error(), "\n")
		messageBuffer.WriteString("\n")
		if len(parts) > 1 {
			for j, part := range parts {
				prefix := "\n "
				if j == 0 {
					prefix = "-"
				}
				messageBuffer.WriteString(fmt.Sprintf("%s %s", prefix, part))
			}
		} else {
			messageBuffer.WriteString(fmt.Sprintf("- %v", err))
		}
	}

	return messageBuffer.String()
}

// Is reports whether any of the internal errors matches the target
func (err ErrAggregate) Is(target error) bool {
	for _, err := range err.errs {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}
