package image

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Image

import (
	"github.com/opspec-io/sdk-golang/expression"
	"github.com/opspec-io/sdk-golang/model"
)

type Image interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallImage *model.SCGContainerCallImage,
		pkgHandle model.PkgHandle,
	) (*model.DCGContainerCallImage, error)
}

func New() Image {
	return _Image{
		expression: expression.New(),
	}
}

type _Image struct {
	expression expression.Expression
}
