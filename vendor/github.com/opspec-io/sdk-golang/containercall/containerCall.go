// Package containercall implements usecases surrounding container calls
package containercall

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ ContainerCall

import (
	"github.com/golang-interfaces/ios"
	"github.com/golang-utils/dircopier"
	"github.com/golang-utils/filecopier"
	"github.com/opspec-io/sdk-golang/interpolater"
	"github.com/opspec-io/sdk-golang/model"
)

type ContainerCall interface {
	// Interpret interprets an SCGContainerCall into a DCGContainerCall
	Interpret(
		scope map[string]*model.Value,
		scgContainerCall *model.SCGContainerCall,
		containerId string,
		rootOpId string,
		pkgPath string,
	) (*model.DCGContainerCall, error)
}

func New(
	rootFSPath string,
) ContainerCall {
	return _ContainerCall{
		dirCopier:    dircopier.New(),
		fileCopier:   filecopier.New(),
		interpolater: interpolater.New(),
		os:           ios.New(),
		rootFSPath:   rootFSPath,
	}
}

type _ContainerCall struct {
	dirCopier    dircopier.DirCopier
	fileCopier   filecopier.FileCopier
	interpolater interpolater.Interpolater
	os           ios.IOS
	rootFSPath   string
}
