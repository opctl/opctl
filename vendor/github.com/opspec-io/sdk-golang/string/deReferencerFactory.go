package string

//go:generate counterfeiter -o ./fakeDeReferencerFactory.go --fake-name fakeDeReferencerFactory ./ deReferencerFactory

import (
	"github.com/opspec-io/sdk-golang/model"
	"github.com/opspec-io/sdk-golang/util/interpolater"
)

type deReferencerFactory interface {
	New(
		scope map[string]*model.Value,
	) interpolater.DeReferencer
}

func newDeReferencerFactory() deReferencerFactory {
	return _deReferencerFactory{}
}

type _deReferencerFactory struct{}

func (drf _deReferencerFactory) New(
	scope map[string]*model.Value,
) interpolater.DeReferencer {
	return newDeReferencer(scope)
}
