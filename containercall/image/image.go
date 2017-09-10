package image

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ Image

import (
	"github.com/opspec-io/sdk-golang/model"
	stringPkg "github.com/opspec-io/sdk-golang/string"
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
		string: stringPkg.New(),
	}
}

type _Image struct {
	string stringPkg.String
}
