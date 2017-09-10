// Package containercall implements usecases surrounding container calls
package containercall

//go:generate counterfeiter -o ./fake.go --fake-name Fake ./ ContainerCall

import (
	"github.com/golang-interfaces/ios"
	"github.com/opspec-io/sdk-golang/containercall/dirs"
	"github.com/opspec-io/sdk-golang/containercall/envvars"
	"github.com/opspec-io/sdk-golang/containercall/files"
	"github.com/opspec-io/sdk-golang/containercall/image"
	"github.com/opspec-io/sdk-golang/containercall/sockets"
	"github.com/opspec-io/sdk-golang/model"
	stringPkg "github.com/opspec-io/sdk-golang/string"
)

type ContainerCall interface {
	// Interpret interprets an SCGContainerCall into a DCGContainerCall
	Interpret(
		scope map[string]*model.Value,
		scgContainerCall *model.SCGContainerCall,
		containerId string,
		rootOpId string,
		pkgHandle model.PkgHandle,
	) (*model.DCGContainerCall, error)
}

func New(
	rootFSPath string,
) ContainerCall {
	return _ContainerCall{
		dirs:       dirs.New(rootFSPath),
		envVars:    envvars.New(),
		files:      files.New(rootFSPath),
		image:      image.New(),
		os:         ios.New(),
		rootFSPath: rootFSPath,
		string:     stringPkg.New(),
		sockets:    sockets.New(),
	}
}

type _ContainerCall struct {
	dirs       dirs.Dirs
	envVars    envvars.EnvVars
	files      files.Files
	image      image.Image
	os         ios.IOS
	rootFSPath string
	string     stringPkg.String
	sockets    sockets.Sockets
}
