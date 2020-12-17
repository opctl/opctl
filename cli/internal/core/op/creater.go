package op

import (
	"path/filepath"

	"github.com/opctl/opctl/sdks/go/opspec"
)

// Creater exposes the "op create" sub command
type Creater interface {
	Create(
		path string,
		description string,
		name string,
	) error
}

// newCreater returns an initialized "op create" sub command
func newCreater() Creater {
	return _creater{}
}

type _creater struct{}

func (ivkr _creater) Create(
	path string,
	description string,
	name string,
) error {
	return opspec.Create(
		filepath.Join(path, name),
		name,
		description,
	)
}
