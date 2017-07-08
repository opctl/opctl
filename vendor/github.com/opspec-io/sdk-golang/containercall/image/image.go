package image

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Image

import (
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type Image interface {
	Interpret(
		scope map[string]*model.Value,
		scgContainerCallImage *model.SCGContainerCallImage,
	) (*model.DCGContainerCallImage, error)
}

func New() Image {
	return _Image{
		interpolater: interpolater.New(),
	}
}

type _Image struct {
	interpolater interpolater.Interpolater
}
